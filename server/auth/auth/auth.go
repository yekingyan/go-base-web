package auth

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	authpb "gService/auth/api/gen/v1"
	"gService/auth/dao"
	"gService/share/id"
	sharetrace "gService/share/trace"
)

// Service implements the proto.AuthServiceServer interface.
type Service struct {
	Logger *zap.Logger
	Mongo  *dao.AuthMongo
	Token  TokenGenerator
}

// TokenGenerator is a token interface.
type TokenGenerator interface {
	GetExpiresIn() int64
	GenerateToken(userID id.UserID) (string, int64, error)
}

// HashPassword turns a plaintext password into a hash.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash checks if a password matches a hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetHashingCost returns the hashing cost used by bcrypt.
func GetHashingCost(hashedPassword []byte) int {
	cost, err := bcrypt.Cost(hashedPassword)
	if err != nil {
		panic(err)
	}
	return cost
}

// Register is a gRPC method.
func (s *Service) Register(ctx context.Context, req *authpb.LoginRequest) (*authpb.RegisterResponse, error) {
	hp, err := HashPassword(req.Password)
	if err != nil {
		s.Logger.Error("Register HashPassword", zap.Error(err))
		return nil, status.Errorf(codes.Unavailable, "unvalid password")
	}
	ok, row, err := s.Mongo.CreateUser(req.Username, hp)
	if !ok || err != nil {
		s.Logger.Error("Register CreateUser", zap.Error(err))
		return nil, status.Errorf(codes.AlreadyExists, "username already exists")
	}

	s.Logger.Info("Register success", zap.Any("row", row))
	return &authpb.RegisterResponse{
		UserId:   row.ID.String(),
		Username: row.Username,
	}, nil
}

// Login  is a gRPC method.
func (s *Service) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	row, err := s.Mongo.GetUserByName(req.Username)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "")
	}
	if ok := CheckPasswordHash(req.Password, row.Password); !ok {
		return nil, status.Errorf(codes.Unauthenticated, "wrong password")
	}

	sessionID, expire, err := s.Token.GenerateToken(row.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "")
	}
	s.Logger.Info("Login success", zap.Any("user_id", row.ID))
	return &authpb.LoginResponse{
		AccessToken: sessionID,
		ExpiresIn:   s.Token.GetExpiresIn(),
		Expire:      expire,
		UserId:      row.ID.String(),
		Username:    row.Username,
	}, nil
}

// Ping is a Server api test.
func (s *Service) Ping(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	s.Logger.Info("ping")
	return &authpb.LoginResponse{
		AccessToken: "access_token ping",
		ExpiresIn:   3600,
	}, nil
}

// UnaryInterceptorWithLog log request
func UnaryInterceptorWithLog(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		logger.Info("UnaryInterceptorWithLog", zap.String("method", info.FullMethod), zap.Any("req", req))
		// get request Header X-Request-ID
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("missing metadata")
		}
		requestID := md.Get(sharetrace.XRequestIDKey)
		if len(requestID) > 0 {
			ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(sharetrace.XRequestIDKey, requestID[0]))
		}
		logger.Info("UnaryInterceptorWithLog", zap.String("method", info.FullMethod), zap.Any("request_id", requestID))
		return handler(ctx, req)
	}
}
