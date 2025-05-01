package main

import (
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/nathabuddhi/ay-com/backend/service-notification/database"
	"github.com/nathabuddhi/ay-com/backend/service-notification/handlers"
	pb "github.com/nathabuddhi/ay-com/backend/service-notification/proto/notification"
	"github.com/nathabuddhi/ay-com/backend/service-notification/rabbitmq"
	"github.com/nathabuddhi/ay-com/backend/service-notification/server"
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
	handler := handlers.NewHandler(db)

	go rabbitmq.InitRabbitMQ(handler)
	notifServer := server.NewNotificationServer(handler)

	s := grpc.NewServer()

	lis, err := net.Listen("tcp", ":5009")
	if err != nil {
		zap.L().Fatal("Failed to listen: " + err.Error())
	}

	pb.RegisterNotificationServiceServer(s, notifServer)
	zap.L().Info("Notification Service gRPC server started successfully.")

	zap.L().Info("Notification Service Running. Listening on port 5009.")
	if err := s.Serve(lis); err != nil {
		zap.L().Fatal("Failed to serve: " + err.Error())
	}
}
