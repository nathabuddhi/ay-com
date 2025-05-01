package handlers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nathabuddhi/ay-com/backend/service-notification/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-notification/proto/notification"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	if db == nil {
		zap.L().Fatal("Database connection is nil!")
	}
	return &Handler{
		DB: db,
	}
}

func (h *Handler) CheckNotificationSettings(user_id string, notif_type string) bool {
	var settings models.NotificationSetting

	err := h.DB.Where("user_id = ?", user_id).First(&settings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			settings = models.NotificationSetting{
				UserId:          user_id,
				NotifLike:       true,
				NotifRepost:     true,
				NotifFollow:     true,
				NotifMention:    true,
				NotifCommunity:  true,
				NotifNewsletter: true,
			}
			h.DB.Create(&settings)
			return true
		} else {
			zap.L().Error("Failed to get notification settings: " + err.Error())
			return false
		}
	}
	if notif_type == "system" {
		return true
	}
	if notif_type == "like" && settings.NotifLike {
		return true
	}
	if notif_type == "community" && settings.NotifCommunity {
		return true
	}
	if notif_type == "follow" && settings.NotifFollow {
		return true
	}
	if notif_type == "repost" && settings.NotifRepost {
		return true
	}
	if notif_type == "mention" && settings.NotifMention {
		return true
	}
	if notif_type == "newsletter" && settings.NotifNewsletter {
		return true
	}
	return false
}

func (h *Handler) CreateNotification(user_id string, title string, content string, from string) {
	generatedId := uuid.New().String()
	for {
		if h.DB.Where("notification_id = ?", generatedId).First(&pb.Notification{}).Error != nil {
			break
		}
		generatedId = uuid.New().String()
	}

	notification := models.Notification{
		NotificationId: generatedId,
		UserId:         user_id,
		Title:          title,
		Content:        content,
		From:           from,
		Read:           false,
		Timestamp:      time.Now(),
	}

	err := h.DB.Create(notification).Error
	if err != nil {
		zap.L().Error("Failed to create notification: " + err.Error())
	} else {
		zap.L().Info("Created notification successfully", zap.String("notification_id", generatedId))
	}
}

func (h *Handler) Notification_GetAllNotifications(ctx context.Context, req *pb.GetAllNotificationsRequest) (*pb.ApiResponseNotification, error) {
	zap.L().Info("User getting all notifications", zap.String("user_id", req.UserId))

	var notifications []models.Notification
	err := h.DB.WithContext(ctx).
		Where("user_id = ?", req.UserId).
		Order("read desc").
		Order("timestamp desc").
		Find(&notifications).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return &pb.ApiResponseNotification{Success: false, Message: "Failed to get notifications."}, nil
	}

	if len(notifications) == 0 {
		return &pb.ApiResponseNotification{Success: true, Message: "You have no notifications yet."}, nil
	}

	allNotificationResponse := &pb.GetAllNotificationResponse{
		Notifications: make([]*pb.Notification, len(notifications)),
	}

	for i, notif := range notifications {
		allNotificationResponse.Notifications[i] = &pb.Notification{
			NotificationId: notif.NotificationId,
			Title:          notif.Title,
			Content:        notif.Content,
			From:           notif.From,
			Read:           notif.Read,
			Timestamp:      notif.Timestamp.String(),
		}
	}

	returnData, err := anypb.New(allNotificationResponse)
	if err != nil {
		return &pb.ApiResponseNotification{
			Success: false,
			Message: "An error occured: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponseNotification{
		Success: true,
		Message: "Get verification requests successful.",
		Data:    returnData,
	}, nil
}

func (h *Handler) Notification_MarkNotificationAsRead(ctx context.Context, req *pb.MarkNotificationAsReadRequest) (*pb.ApiResponseNotification, error) {
	zap.L().Info("User " + req.UserId + " marking notification " + req.NotificationId + " as read.")

	var notification models.Notification
	err := h.DB.WithContext(ctx).
		Where("notification_id = ?", req.NotificationId).
		Where("user_id = ?", req.UserId).
		First(&notification).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.ApiResponseNotification{Success: false, Message: "Notification not found."}, nil
		}
		zap.L().Error("Failed to mark notification as read: " + err.Error())
		return &pb.ApiResponseNotification{Success: false, Message: "Failed to mark notification as read."}, nil
	}

	if notification.Read {
		return &pb.ApiResponseNotification{Success: false, Message: "Notification already marked as read."}, nil
	}

	notification.Read = true
	err = h.DB.Save(&notification).Error
	if err != nil {
		zap.L().Error("Failed to mark notification as read: " + err.Error())
		return &pb.ApiResponseNotification{Success: false, Message: "Failed to mark notification as read."}, nil
	}

	return &pb.ApiResponseNotification{Success: true, Message: "Successfully marked notification as read."}, nil
}

func (h *Handler) Notification_DeleteNotification(ctx context.Context, req *pb.DeleteNotificationRequest) (*pb.ApiResponseNotification, error) {
	zap.L().Info("User " + req.UserId + " deleting notification " + req.NotificationId)

	var notification models.Notification
	err := h.DB.WithContext(ctx).
		Where("notification_id = ?", req.NotificationId).
		Where("user_id = ?", req.UserId).
		First(&notification).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.ApiResponseNotification{Success: false, Message: "Notification not found."}, nil
		}
		zap.L().Error("Failed to mark notification as read: " + err.Error())
		return &pb.ApiResponseNotification{Success: false, Message: "Failed to mark notification as read."}, nil
	}

	err = h.DB.Delete(&notification).Error
	if err != nil {
		zap.L().Error("Failed to delete notification: " + err.Error())
		return &pb.ApiResponseNotification{Success: false, Message: "Failed to delete notification."}, nil
	}
	return &pb.ApiResponseNotification{Success: true, Message: "Successfully deleted notification."}, nil
}

func (h *Handler) Notification_ClearNotifications(ctx context.Context, req *pb.ClearNotificationsRequest) (*pb.ApiResponseNotification, error) {
	zap.L().Info("User " + req.UserId + " is clearing notifications.")

	err := h.DB.WithContext(ctx).
		Where("user_id = ?", req.UserId).
		Delete(&models.Notification{}).Error

	if err != nil {
		zap.L().Error("Failed to clear notifications for user " + req.UserId + ": " + err.Error())
		return &pb.ApiResponseNotification{Success: false, Message: "Failed to clear notifications."}, nil
	}

	return &pb.ApiResponseNotification{Success: true, Message: "Successfully cleared all notifications."}, nil
}
