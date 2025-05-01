package handlers

import (
	"net/http"
	"sync"

	"github.com/nathabuddhi/ay-com/backend/api-gateway/middleware"
	pb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/notification"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
