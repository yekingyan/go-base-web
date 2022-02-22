package main

import (
	"context"
	"net"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	authpb "gService/auth/api/gen/v1"
	"gService/auth/auth"
	"gService/auth/dao"
	stoken "gService/auth/token"
	sharelog "gService/share/log"
)

// PORT is the port to listen on.
const PORT string = ":9001"

// MongoURL is the mongo url.
const MongoURL string = "mongodb://localhost:27017"

// DatabaseName is the name of the database.
const DatabaseName string = "gservice"

// Logger is the global logger.
var Logger *zap.Logger

// TokenExpire is the token expire time.
const TokenExpire int64 = 60 * 60 * 2

func main() {
	Logger := sharelog.InitZapLog("auth")
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		Logger.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(auth.UnaryInterceptorWithLog(Logger)),
	)

	ctx := context.Background()
	mc, err := mongo.Connect(ctx,
		options.Client().ApplyURI(MongoURL))
	if err != nil {
		Logger.Fatal("failed to connect to mongo:", zap.Error(err))
	}

	am := dao.NewMongo(mc.Database(DatabaseName), Logger)
	authpb.RegisterAuthServiceServer(s, &auth.Service{
		Logger: Logger,
		Mongo:  am,
		Token: stoken.NewSessionToken(
			TokenExpire,
			am,
			Logger,
		),
	})

	Logger.Info("starting auth service", zap.String("port", PORT))
	err = s.Serve(lis)
	Logger.Fatal("failed to serve", zap.Error(err))
}
