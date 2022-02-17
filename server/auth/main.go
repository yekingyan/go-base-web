package main

import (
	"context"
	"log"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	authpb "gService/auth/api/gen/v1"
	"gService/auth/auth"
)

// PORT is the port to listen on.
const PORT string = ":9001"

func filter(ctx context.Context,
	req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	log.Println("fileter:", info)
	return handler(ctx, req)
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(filter),
	)
	authpb.RegisterAuthServiceServer(s, &auth.Service{
		Logger: logger,
	})

	logger.Info("starting auth service", zap.String("port", PORT))
	err = s.Serve(lis)
	logger.Fatal("failed to serve", zap.Error(err))
}
