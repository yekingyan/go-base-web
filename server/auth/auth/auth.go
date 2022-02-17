package auth

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	authpb "gService/auth/api/gen/v1"
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
	return &authpb.LoginResponse{}, nil
}

// Interceptor log request
func Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	fmt.Println("interceptor:", zap.String("request", fmt.Sprintf("%+v", req)))
	return handler(ctx, req)
}
