package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	redis_client "github.com/nathabuddhi/ay-com/backend/api-gateway/redis"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/types"
	"go.uber.org/zap"
)

var (
	USER_SERVICE_PATH         string
	NOTIFICATION_SERVICE_PATH string
	FLASK_SERVICE_PATH        string
	MEDIA_SERVICE_PATH        string
	THREAD_SERVICE_PATH       string
)

func InitEnvironmentVariables() {
	err := godotenv.Load()
	if err != nil {
		zap.L().Fatal("Error loading .env file")
	}

	USER_SERVICE_PATH = os.Getenv("USER_SERVICE_PATH")
	FLASK_SERVICE_PATH = os.Getenv("FLASK_SERVICE_PATH")
	NOTIFICATION_SERVICE_PATH = os.Getenv("NOTIFICATION_SERVICE_PATH")
	MEDIA_SERVICE_PATH = os.Getenv("MEDIA_SERVICE_PATH")
	THREAD_SERVICE_PATH = os.Getenv("THREAD_SERVICE_PATH")

	zap.L().Info("Service Paths Loaded Successfully.")
}

func returnErrorResponse(w http.ResponseWriter, message string) {
	response := types.ApiResponse{
		Success: false,
		Message: message,
		Payload: nil,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}

func checkRedisData(key string, w http.ResponseWriter) bool {
	redisProfile := redis_client.GetCache(key)

	if redisProfile != nil {
		response := types.ApiResponse{
			Success: true,
			Message: "Get Cached Data successful.",
			Payload: redisProfile,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return true
	}

	return false
}

func createContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	return ctx, cancel
}
