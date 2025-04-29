package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/nathabuddhi/ay-com/backend/util-redis/rabbitmq"
	"github.com/nathabuddhi/ay-com/backend/util-redis/redis_client"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		zap.L().Fatal("Failed to connect to RabbitMQ: " + err.Error())
	}
	defer conn.Close()

	rabbitmq.InitSetRedisChannel(conn)
	rabbitmq.InitDeleteRedisChannel(conn)

	rabbitmq.StartConsumingSet()
	rabbitmq.StartConsumingDelete()

	zap.L().Info("Redis Service Running. Listening for RabbitMQ messages.")

	<-rabbitmq.Forever
}
