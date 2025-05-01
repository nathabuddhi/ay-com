package handlers

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-notification/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-notification/proto/notification"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
)

func (h *Handler) Notification_GetSettings(ctx context.Context, req *pb.GetNotificationSettingsRequest) (*pb.ApiResponseNotification, error) {
	zap.L().Info("User " + req.UserId + " is getting notification settings.")

	var notificationSetting models.NotificationSetting
	if err := h.DB.Where("user_id = ?", req.UserId).First(&notificationSetting).Error; err != nil {
		return &pb.ApiResponseNotification{
			Success: false,
			Message: "An unknown error occured.",
			Data:    nil,
		}, nil
	}

	notificationSettingsReturn := &pb.NotificationSettings{
		NotifLike:       notificationSetting.NotifLike,
		NotifRepost:     notificationSetting.NotifRepost,
		NotifFollow:     notificationSetting.NotifFollow,
		NotifMention:    notificationSetting.NotifMention,
		NotifCommunity:  notificationSetting.NotifCommunity,
		NotifNewsletter: notificationSetting.NotifNewsletter,
	}

	returnData, err := anypb.New(notificationSettingsReturn)
	if err != nil {
		return &pb.ApiResponseNotification{
			Success: false,
			Message: "An error occured: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponseNotification{
		Success: true,
		Message: "Get Notification Settings successful.",
		Data:    returnData,
	}, nil
}

func (h *Handler) Notification_UpdateSettings(ctx context.Context, req *pb.UpdateNotificationSettingsRequest) (*pb.ApiResponseNotification, error) {
	zap.L().Info("User " + req.UserId + " is getting notification settings.")

	var notificationSetting models.NotificationSetting
	if err := h.DB.Where("user_id = ?", req.UserId).First(&notificationSetting).Error; err != nil {
		return &pb.ApiResponseNotification{
			Success: false,
			Message: "An unknown error occured.",
			Data:    nil,
		}, nil
	}

	notificationSetting.NotifLike = req.NotifLike
	notificationSetting.NotifRepost = req.NotifRepost
	notificationSetting.NotifFollow = req.NotifFollow
	notificationSetting.NotifMention = req.NotifMention
	notificationSetting.NotifCommunity = req.NotifCommunity
	notificationSetting.NotifNewsletter = req.NotifNewsletter

	if err := h.DB.Save(&notificationSetting).Error; err != nil {
		zap.L().Error("Failed to update notification settings: " + err.Error())
		return &pb.ApiResponseNotification{
			Success: false,
			Message: "An unknown error occured.",
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponseNotification{
		Success: true,
		Message: "Update Notificatiion Settings successful.",
		Data:    nil,
	}, nil
}
