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
	unionpb "gService/union/api/gen/v1"
)

// PORT is the getway port.
const PORT string = ":7000"

// Logger is the global logger.
var Logger *zap.Logger

var serverConfig = map[string] struct {
	addr     string
	register func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
}{
	"auth": {
		addr:     "localhost:7001",
		register: authpb.RegisterAuthServiceHandlerFromEndpoint,
	},
	"union": {
		addr:     "localhost:7002",
		register: unionpb.RegisterUnionServiceHandlerFromEndpoint,
	},
}

func tracingWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(sharetrace.XRequestIDKey)
		if requestID == "" {
			requestID = sharetrace.NewRequestID()
		}
		r.Header.Set(sharetrace.XRequestIDKey, requestID)
		Logger.Info("request", zap.String("method", r.Method), zap.String("url", r.URL.String()), zap.String("request_id", requestID))
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

	for name, s := range serverConfig {
		err := s.register(ctx, mux, s.addr, []grpc.DialOption{grpc.WithInsecure()})
		if err != nil {
			Logger.Fatal("failed to register service:", zap.String("name", name), zap.Error(err))
		}
		Logger.Info("registered service:", zap.String("name", name))
	}

	fmt.Println("gateway listening on", PORT)
	Logger.Sugar().Fatal(http.ListenAndServe(PORT, tracingWrapper(mux)))
}
