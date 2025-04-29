package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/handlers"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/middleware"
	redis_client "github.com/nathabuddhi/ay-com/backend/api-gateway/redis"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/supabase"
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
	supabase.InitSupabase()

	r := mux.NewRouter()

	r.HandleFunc("/user/login", handlers.User_Login).Methods("POST")
	r.HandleFunc("/user/register", handlers.User_Register).Methods("POST")
	r.HandleFunc("/user/requestverificationcode", handlers.User_RequestVerificationCode).Methods("POST")
	r.HandleFunc("/user/validateverificationcode", handlers.User_ValidateVerificationCode).Methods("POST")

	secured := r.PathPrefix("/").Subrouter()
	secured.Use(middleware.JwtAuthMiddleware)
	// secured.HandleFunc("/user/getprofile/:id", handlers.User_GetProfile).Methods("GET")

	httpHandler := allowCors(r)

	zap.L().Info("API Gateway Running. Listening on port 5000.")

	if err := http.ListenAndServe(":"+"5000", httpHandler); err != nil {
		zap.L().Fatal("Server failed to start", zap.Error(err))
	}
}

func allowCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		h.ServeHTTP(w, r)
	})
}
