package main

import (
	"context"
	"crypto/rsa"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
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

// PrivateKeyPath is the path of private key.
const PrivateKeyPath = "auth/private.key"

func getPk(logger *zap.Logger) *rsa.PrivateKey {
	pkFile, err := os.Open(PrivateKeyPath)
	if err != nil {
		logger.Fatal("cannot open private key", zap.Error(err))
		panic(err)
	}
	pkBytes, err := ioutil.ReadAll(pkFile)
	if err != nil {
		logger.Fatal("cannot read private key", zap.Error(err))
		panic(err)
	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
	if err != nil {
		logger.Fatal("cannot parse private key", zap.Error(err))
		panic(err)
	}
	return privKey
}

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
		Token: stoken.NewJWTToken(
			TokenExpire,
			time.Now,
			getPk(Logger),
		),
	})

	Logger.Info("starting auth service", zap.String("port", PORT))
	err = s.Serve(lis)
	Logger.Fatal("failed to serve", zap.Error(err))
}
