package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	pb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/user"
	redis_client "github.com/nathabuddhi/ay-com/backend/api-gateway/redis"
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

func processResponseWithoutPayload(resp *pb.ApiResponse, err error, w http.ResponseWriter) {
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	response := &types.ApiResponse{
		Success: resp.Success,
		Message: resp.Message,
		Payload: nil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func processResponseWithPayload[T any](resp *pb.ApiResponse, err error, w http.ResponseWriter) {
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	decodedObject := new(T)
	if _, ok := any(decodedObject).(proto.Message); !ok {
		zap.L().Error("Type does not implement proto.Message", zap.String("type", fmt.Sprintf("%T", decodedObject)))
		returnErrorResponse(w, "Failed to decode response data.")
		return
	}

	if err := proto.Unmarshal(resp.Data.GetValue(), any(decodedObject).(proto.Message)); err != nil {
		zap.L().Error("Failed to unmarshal protobuf", zap.Error(err))
		returnErrorResponse(w, "Failed to decode response data.")
		return
	}

	response := &types.ApiResponse{
		Success: resp.Success,
		Message: resp.Message,
		Payload: decodedObject,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func processRequest[T any](r *http.Request, w http.ResponseWriter) (resp *T, client pb.UserServiceClient) {
	conn := getUserServiceConn()
	client = pb.NewUserServiceClient(conn)

	var req T
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to encode API request", zap.Error(err))
		returnErrorResponse(w, "Failed to encode API request: "+err.Error())
	}
	return &req, client
}

func checkRedisData(key string, w http.ResponseWriter) bool {
	redisProfile := redis_client.GetCache(key)

	if redisProfile != nil {
		response := types.ApiResponse{
			Success: true,
			Message: "Get User Profile successful.",
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
