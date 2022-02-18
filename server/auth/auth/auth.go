package auth

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	authpb "gService/auth/api/gen/v1"
	sharetrace "gService/share/trace"
)

// Service implements the proto.AuthServiceServer interface.
type Service struct {
	Logger *zap.Logger
}

// Login in auth service.
func (s *Service) Login(context.Context, *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	s.Logger.Info("Login")
	return &authpb.LoginResponse{
		AccessToken: "access_token",
		ExpiresIn:   3600,
	}, nil
}

// Ping is a Server api test.
func (s *Service) Ping(context.Context, *authpb.LoginRequest) (*authpb.LoginResponse, error) {
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
