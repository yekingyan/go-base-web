package main

import (
	"gService/share/auth"
	sharelog "gService/share/log"
	unionpb "gService/union/api/gen/v1"
	"gService/union/union"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// PORT is the port to listen on.
const PORT string = ":7002"

// Logger is the global logger.
var Logger *zap.Logger

// PublicKeyPath is the path of public key.
const PublicKeyPath = "../share/auth/public.key"

func main() {
	Logger := sharelog.InitZapLog("union")
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		Logger.Fatal("failed to listen", zap.Error(err))
	}

	in, err := auth.GetInterceptor(PublicKeyPath)
	if err != nil {
		Logger.Fatal("failed to get auth interceptor", zap.Error(err))
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(in),
	)

	// var _ unionpb.UnionServiceServer = (*union.Service)(nil)
	unionpb.RegisterUnionServiceServer(s, &union.Service{
		Logger: Logger,
	})

	Logger.Info("starting union service", zap.String("port", PORT))
	err = s.Serve(lis)
	Logger.Fatal("failed to serve", zap.Error(err))
}
