package main

import (
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/nathabuddhi/ay-com/backend/service-user/database"
	emailpb "github.com/nathabuddhi/ay-com/backend/service-user/proto/email"
	redispb "github.com/nathabuddhi/ay-com/backend/service-user/proto/redis"
	pb "github.com/nathabuddhi/ay-com/backend/service-user/proto/user"
	"github.com/nathabuddhi/ay-com/backend/service-user/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	db := database.InitDB()

	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		zap.L().Fatal("Failed to listen: " + err.Error())
	}

	redisConn, err := grpc.NewClient(os.Getenv("REDIS_SERVICE_PATH"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Fatal("Failed to connect to Redis service: " + err.Error())
	}
	redisClient := redispb.NewRedisServiceClient(redisConn)
	zap.L().Info("Connected to Redis service.")

	emailConn, err := grpc.NewClient(os.Getenv("EMAIL_SERVICE_PATH"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Fatal("Failed to connect to Email service: " + err.Error())
	}
	emailClient := emailpb.NewEmailServiceClient(emailConn)
	zap.L().Info("Connected to Email service.")

	userServer := server.NewUserServer(db, redisClient, emailClient)

	s := grpc.NewServer()

	pb.RegisterUserServiceServer(s, userServer)
	zap.L().Info("User Service gRPC server started successfully.")

	zap.L().Info("User Service Running. Listening on port 5001.")
	if err := s.Serve(lis); err != nil {
		zap.L().Fatal("Failed to serve: " + err.Error())
	}
}
