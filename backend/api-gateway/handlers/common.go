package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/types"
	"go.uber.org/zap"
)

var (
	USER_SERVICE_PATH  string
	FLASK_SERVICE_PATH string
)

func InitEnvironmentVariables() {
	err := godotenv.Load()
	if err != nil {
		zap.L().Fatal("Error loading .env file")
	}

	USER_SERVICE_PATH = os.Getenv("USER_SERVICE_PATH")
	FLASK_SERVICE_PATH = os.Getenv("FLASK_SERVICE_PATH")

	zap.L().Info("Service Paths Loaded Successfully.")
}


func returnErrorResponse(w http.ResponseWriter, message string) {
	response := types.ApiResponse{
		Success: false,
		Message: message,
		Payload: nil,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(response)
}
