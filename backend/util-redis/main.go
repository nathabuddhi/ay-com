package main

import (
	"net"
	"os"

	"github.com/joho/godotenv"
	pb "github.com/nathabuddhi/ay-com/backend/util-redis/proto"
	"github.com/nathabuddhi/ay-com/backend/util-redis/redis_client"
	"github.com/nathabuddhi/ay-com/backend/util-redis/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

func getLogLevel() zapcore.Level {
	levelStr := os.Getenv("APP_LOG_LEVEL")
	switch levelStr {
	case "info":
		return zapcore.InfoLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func initLogger() {
	config := zap.Config{
		Level:         zap.NewAtomicLevelAt(getLogLevel()),
		Development:   false,
		Encoding:      "json",
		OutputPaths:   []string{"stdout"},
		EncoderConfig: zap.NewProductionEncoderConfig(),
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}

func main() {
	initLogger()
	defer zap.L().Sync()

	err := godotenv.Load()
	if err != nil {
		zap.L().Error("Error loading .env file")
	}

	redis_client.InitRedis()

	lis, err := net.Listen("tcp", ":5012")
	if err != nil {
		zap.L().Fatal("Failed to listen: " + err.Error())
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRedisServiceServer(grpcServer, server.NewRedisServer())

	zap.L().Info("Redis Service gRPC server started successfully.")

	zap.L().Info("Redis Service Running. Listening on port 5012.")

	if err := grpcServer.Serve(lis); err != nil {
		zap.L().Fatal("Failed to serve: " + err.Error())
	}
}
