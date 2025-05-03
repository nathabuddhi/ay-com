package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/handlers"
	redis_client "github.com/nathabuddhi/ay-com/backend/api-gateway/redis"
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

	handlers.InitEnvironmentVariables()
	redis_client.InitRedis()

	r := handlers.InitRoutes()

	httpHandler := allowCors(r)

	zap.L().Info("API Gateway Running. Listening on port 5000.")

	if err := http.ListenAndServe(":"+"5000", httpHandler); err != nil {
		zap.L().Fatal("Server failed to start", zap.Error(err))
	}
}

func allowCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
