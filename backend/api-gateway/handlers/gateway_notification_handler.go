package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/nathabuddhi/ay-com/backend/api-gateway/middleware"
	pb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/notification"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/types"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

var (
	notifServiceConn *grpc.ClientConn
	notifServiceOnce sync.Once
)

func getNotifServiceConn() *grpc.ClientConn {
	notifServiceOnce.Do(func() {
		var err error
		notifServiceConn, err = grpc.NewClient(NOTIFICATION_SERVICE_PATH, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			zap.L().Fatal("Failed to connect to notification service: " + err.Error())
		}
	})
	return notifServiceConn
}

func processNotificationResponseWithoutPayload(resp *pb.ApiResponseNotification, err error, w http.ResponseWriter) {
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

func processNotificationResponseWithPayload[T any](resp *pb.ApiResponseNotification, err error, w http.ResponseWriter) {
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

func processNotificationRequest[T any](r *http.Request, w http.ResponseWriter) (resp *T, client pb.NotificationServiceClient) {
	conn := getNotifServiceConn()
	client = pb.NewNotificationServiceClient(conn)

	var req T
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to encode API request", zap.Error(err))
		returnErrorResponse(w, "Failed to encode API request: "+err.Error())
	}
	return &req, client
}

func Notification_GetSettings(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Notification Get Settings is called.")

	req, client := processNotificationRequest[pb.GetNotificationSettingsRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Notification_GetSettings(ctx, req)

	processNotificationResponseWithPayload[pb.NotificationSettings](resp, err, w)
}

func Notification_UpdateSettings(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Notification Update Settings is called.")

	req, client := processNotificationRequest[pb.UpdateNotificationSettingsRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Notification_UpdateSettings(ctx, req)

	processNotificationResponseWithoutPayload(resp, err, w)
}

func Notification_GetAllNotifications(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Notification Update Settings is called.")

	req, client := processNotificationRequest[pb.GetAllNotificationsRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Notification_GetAllNotifications(ctx, req)

	processNotificationResponseWithPayload[pb.GetAllNotificationResponse](resp, err, w)
}

func Notification_ClearNotifications(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Notification Clear Notifications is called.")

	req, client := processNotificationRequest[pb.ClearNotificationsRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Notification_ClearNotifications(ctx, req)

	processNotificationResponseWithoutPayload(resp, err, w)
}

func Notification_DeleteNotification(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Notification Delete Notification is called.")

	req, client := processNotificationRequest[pb.DeleteNotificationRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Notification_DeleteNotification(ctx, req)

	processNotificationResponseWithoutPayload(resp, err, w)
}

func Notification_MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Notification Mark Notification as Read is called.")

	req, client := processNotificationRequest[pb.MarkNotificationAsReadRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Notification_MarkNotificationAsRead(ctx, req)

	processNotificationResponseWithoutPayload(resp, err, w)
}
