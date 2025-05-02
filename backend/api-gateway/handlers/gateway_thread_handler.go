package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/nathabuddhi/ay-com/backend/api-gateway/middleware"
	pb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/thread"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/types"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

var (
	threadServiceConn *grpc.ClientConn
	threadServiceOnce sync.Once
)

func getThreadServiceConn() *grpc.ClientConn {
	threadServiceOnce.Do(func() {
		var err error
		threadServiceConn, err = grpc.NewClient(THREAD_SERVICE_PATH, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			zap.L().Fatal("Failed to connect to thread service: " + err.Error())
		}
	})
	return threadServiceConn
}

func processThreadResponseWithoutPayload(resp *pb.ApiResponseThread, err error, w http.ResponseWriter) {
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

func processThreadResponseWithPayload[T any](resp *pb.ApiResponseThread, err error, w http.ResponseWriter) {
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

func processThreadRequest[T any](r *http.Request, w http.ResponseWriter) (resp *T, client pb.ThreadServiceClient) {
	conn := getThreadServiceConn()
	client = pb.NewThreadServiceClient(conn)

	var req T
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to encode API request", zap.Error(err))
		returnErrorResponse(w, "Failed to encode API request: "+err.Error())
	}
	return &req, client
}

func Thread_GetAllThreads(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread GetAllThreads is called.")

	req, client := processThreadRequest[pb.GetAllThreadsRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_GetAllThreads(ctx, req)

	processThreadResponseWithPayload[pb.GetThreadsResponse](resp, err, w)
}
