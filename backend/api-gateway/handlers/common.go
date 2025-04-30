package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	pb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/user"
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

func forwardRequest[TReq any, TRes any](w http.ResponseWriter, req *TReq, grpcCall func(context.Context, *TReq) (*pb.ApiResponse, error)) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res, err := grpcCall(ctx, req)
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	var payload any
	if res.Data != nil {
		newResponse := new(TRes)
		if pm, ok := any(newResponse).(proto.Message); ok {
			err := res.Data.UnmarshalTo(pm)
			if err != nil {
				zap.L().Error("Failed to unmarshal response data", zap.Error(err))
				returnErrorResponse(w, "Failed to process response data")
				return
			}
			switch v := any(newResponse).(type) {
			case **pb.String:
				payload = (*v).Value
			case **pb.UserProfile:
				payload = *v // Use the UserProfile struct directly as the payload
			default:
				payload = newResponse
			}
		} else {
			payload = nil
		}
	}

	response := types.ApiResponse{
		Success: res.Success,
		Message: res.Message,
		Payload: payload,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
