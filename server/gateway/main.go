package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"

	authpb "gService/auth/api/gen/v1"
	sharelog "gService/share/log"
	sharetrace "gService/share/trace"
)

// PORT is the getway port.
const PORT string = ":9000"

// AuthPoint is the auth service address.
const AuthPoint string = "localhost:9001"

// Logger is the global logger.
var Logger *zap.Logger

func tracingWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(sharetrace.XRequestIDKey)
		if requestID == "" {
			requestID = sharetrace.NewRequestID()
		}
		r.Header.Set(sharetrace.XRequestIDKey, requestID)
		Logger.Sugar().Info("request", zap.String("method", r.Method), zap.String("url", r.URL.String()), zap.String("request_id", requestID))
		h.ServeHTTP(w, r)
	})
}

func requestIDAnnotator(ctx context.Context, r *http.Request) metadata.MD {
	requestID := r.Header.Get(sharetrace.XRequestIDKey)
	if requestID == "" {
		requestID = sharetrace.NewRequestID()
	}
	return metadata.Pairs(sharetrace.XRequestIDKey, requestID)
}

func main() {
	Logger = sharelog.InitZapLog("gateway")
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
		),
		runtime.WithMetadata(requestIDAnnotator),
	)
	err := authpb.RegisterAuthServiceHandlerFromEndpoint(
		ctx, mux, AuthPoint,
		[]grpc.DialOption{
			grpc.WithInsecure(), // Ignore certificate errors
			// grpc.WithChainUnaryInterceptor(),
		},
	)
	if err != nil {
		Logger.Sugar().Fatal("failed to register auth service:", err)
	} else {
		Logger.Sugar().Info("connected to auth service:", AuthPoint)
	}
	fmt.Println("gateway listening on", PORT)
	Logger.Sugar().Fatal(http.ListenAndServe(PORT, tracingWrapper(mux)))
}
