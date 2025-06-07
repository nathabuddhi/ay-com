package main

import (
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/nathabuddhi/ay-com/backend/service-user/database"
	"github.com/nathabuddhi/ay-com/backend/service-user/handlers"
	pb "github.com/nathabuddhi/ay-com/backend/service-user/proto/user"
	rabbitmqreceive "github.com/nathabuddhi/ay-com/backend/service-user/rabbitmqreceive"
	rabbitmq "github.com/nathabuddhi/ay-com/backend/service-user/rabbitmqsend"
	"github.com/nathabuddhi/ay-com/backend/service-user/server"
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

	db := database.InitDB()
	handler := handlers.NewHandlers(db)
	rabbitmq.InitRabbitMQ()

	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		zap.L().Fatal("Failed to listen: " + err.Error())
	}

	userServer := server.NewUserServer(handler)

	s := grpc.NewServer()

	pb.RegisterUserServiceServer(s, userServer)
	zap.L().Info("User Service gRPC server started successfully.")

	zap.L().Info("User Service Running. Listening on port 5001.")
	go rabbitmqreceive.InitHandler(handler)
	if err := s.Serve(lis); err != nil {
		zap.L().Fatal("Failed to serve: " + err.Error())
	}
}
