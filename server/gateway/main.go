package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	authpb "gService/auth/api/gen/v1"
)

// PORT is the getway port.
const PORT string = ":9000"

// AuthPoint is the auth service address.
const AuthPoint string = "localhost:9001"

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			&runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseEnumNumbers: true,
				},
			},
		))
	err = authpb.RegisterAuthServiceHandlerFromEndpoint(
		ctx, mux, AuthPoint,
		[]grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		logger.Sugar().Fatal("failed to register auth service:", zap.Error(err))
	} else {
		logger.Sugar().Info("connected to auth service:", AuthPoint)
	}
	fmt.Println("gateway listening on", PORT)
	logger.Sugar().Fatal(http.ListenAndServe(PORT, mux))
}
