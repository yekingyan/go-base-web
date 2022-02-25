package auth

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	sharetoken "gService/share/auth/token"
)

const authorizationHeader = "authorization"
const authPerfix = "Bearer "

// TokenVerifier verifies a token.
type TokenVerifier interface {
	Verify(token string) (string, error)
}

type interceptor struct {
	verifier TokenVerifier
}

// GetInterceptor returns a new auth interceptor.
func GetInterceptor(publicKeyPath string) (grpc.UnaryServerInterceptor, error) {
	f, err := os.Open(publicKeyPath)
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("cannot open public key: %v", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("cannot read public key: %v", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(b)
	if err != nil {
		return nil, fmt.Errorf("cannot parse public key: %v", err)
	}

	i := &interceptor{
		verifier: &sharetoken.JWTTokenVerifier{
			PublicKey: publicKey,
		}}
	return i.VerifyRequst, nil
}

func getTokenFromCtx(ctx context.Context) (string, error) {
	unauthenticated := status.Error(codes.Unauthenticated, "")
	// 从ctx中的metadata(request.HEADER)中获取token
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", unauthenticated
	}
	var tkn string
	// 可以有多个authorization Header
	for _, v := range md["authorization"] {
		if strings.HasPrefix(v, "Bearer ") {
			tkn = v[len("Bearer "):]
			break
		}
	}
	if tkn == "" {
		return "", unauthenticated
	}
	return tkn, nil
}

func (i *interceptor) VerifyRequst(ctx context.Context,
	req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	tkn, err := getTokenFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	// jwt验证
	userID, err := i.verifier.Verify(tkn)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "")
	}
	// 注入userID
	ctx = ContextWithUserID(ctx, userID)
	return handler(ctx, req)
}

type userIDKey struct{}

// ContextWithUserID returns a new context with userID.
func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

// UserIDFromContext returns userID from context.
func UserIDFromContext(ctx context.Context) (userID string, err error) {
	userID, ok := ctx.Value(userIDKey{}).(string)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "")
	}
	return userID, nil
}
