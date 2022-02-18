package main

import (
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	authpb "gService/auth/api/gen/v1"
	"gService/auth/auth"
	sharelog "gService/share/log"
)

// PORT is the port to listen on.
const PORT string = ":9001"

// Logger is the global logger.
var Logger *zap.Logger

func main() {
	Logger := sharelog.InitZapLog("auth")
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		Logger.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(auth.UnaryInterceptorWithLog(Logger)),
	)
	authpb.RegisterAuthServiceServer(s, &auth.Service{
		Logger: Logger,
	})

	Logger.Info("starting auth service", zap.String("port", PORT))
	err = s.Serve(lis)
	Logger.Fatal("failed to serve", zap.Error(err))
}
