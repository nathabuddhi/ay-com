package server

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-notification/handlers"
	pb "github.com/nathabuddhi/ay-com/backend/service-notification/proto/notification"
)

type NotificationServer struct {
	pb.UnimplementedNotificationServiceServer
	Handlers *handlers.Handler
}

func NewNotificationServer(handler *handlers.Handler) *NotificationServer {
	return &NotificationServer{
		Handlers: handler,
	}
}

func (s *NotificationServer) Notification_GetAllNotifications(ctx context.Context, req *pb.GetAllNotificationsRequest) (*pb.ApiResponseNotification, error) {
	return s.Handlers.Notification_GetAllNotifications(ctx, req)
}

func (s *NotificationServer) Notification_ClearNotifications(ctx context.Context, req *pb.ClearNotificationsRequest) (*pb.ApiResponseNotification, error) {
	return s.Handlers.Notification_ClearNotifications(ctx, req)
}

func (s *NotificationServer) Notification_MarkNotificationAsRead(ctx context.Context, req *pb.MarkNotificationAsReadRequest) (*pb.ApiResponseNotification, error) {
	return s.Handlers.Notification_MarkNotificationAsRead(ctx, req)
}

func (s *NotificationServer) Notification_DeleteNotification(ctx context.Context, req *pb.DeleteNotificationRequest) (*pb.ApiResponseNotification, error) {
	return s.Handlers.Notification_DeleteNotification(ctx, req)
}

func (s *NotificationServer) Notification_GetSettings(ctx context.Context, req *pb.GetNotificationSettingsRequest) (*pb.ApiResponseNotification, error) {
	return s.Handlers.Notification_GetSettings(ctx, req)
}

func (s *NotificationServer) Notification_UpdateSettings(ctx context.Context, req *pb.UpdateNotificationSettingsRequest) (*pb.ApiResponseNotification, error) {
	return s.Handlers.Notification_UpdateSettings(ctx, req)
}
