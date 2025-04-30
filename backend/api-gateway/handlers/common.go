package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/types"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
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

func decodeResponse[T any](binaryData []byte) (*T, error) {
	decodedObject := new(T)

	if _, ok := any(decodedObject).(proto.Message); !ok {
		zap.L().Error("Type does not implement proto.Message", zap.String("type", fmt.Sprintf("%T", decodedObject)))
		return nil, fmt.Errorf("type %T does not implement proto.Message", decodedObject)
	}

	// Unmarshal the binary data into decodedObject
	if err := proto.Unmarshal(binaryData, any(decodedObject).(proto.Message)); err != nil {
		zap.L().Error("Failed to unmarshal protobuf", zap.Error(err))
		return nil, err
	}

	return decodedObject, nil
}
