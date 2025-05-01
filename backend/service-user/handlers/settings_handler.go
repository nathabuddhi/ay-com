package handlers

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-user/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-user/proto/user"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
)

func (h *Handlers) User_GetSettings(ctx context.Context, req *pb.GetSettingsRequest) (*pb.ApiResponse, error) {
	zap.L().Info("User " + req.UserId + " is getting settings")

	var userSettings models.UserSetting
	err := h.DB.WithContext(ctx).Where("user_id = ?", req.UserId).First(&userSettings).Error
	if err != nil && err.Error() == "record not found" {
		userSettings = models.UserSetting{
			UserId:          req.UserId,
			FontSize:        "medium",
			FontColor:       "black",
			NotifLike:       true,
			NotifRepost:     true,
			NotifFollow:     true,
			NotifMention:    true,
			NotifCommunity:  true,
			NotifNewsletter: true,
		}

		err = h.DB.WithContext(ctx).Create(&userSettings).Error
		if err != nil {
			return &pb.ApiResponse{
				Success: false,
				Message: "An error occurred while creating user settings: " + err.Error(),
				Data:    nil,
			}, nil
		}

		userSettingsReturn := &pb.UserSettings{
			FontSize:        "medium",
			FontColor:       "black",
			NotifLike:       true,
			NotifRepost:     true,
			NotifFollow:     true,
			NotifMention:    true,
			NotifCommunity:  true,
			NotifNewsletter: true,
		}

		returnData, err := anypb.New(userSettingsReturn)
		if err != nil {
			return &pb.ApiResponse{
				Success: false,
				Message: "An error occured: " + err.Error(),
				Data:    nil,
			}, nil
		}

		return &pb.ApiResponse{
			Success: true,
			Message: "Get User Settings successful.",
			Data:    returnData,
		}, nil
	}

	userSettingsReturn := &pb.UserSettings{
		FontSize:        userSettings.FontSize,
		FontColor:       userSettings.FontColor,
		NotifLike:       userSettings.NotifLike,
		NotifRepost:     userSettings.NotifRepost,
		NotifFollow:     userSettings.NotifFollow,
		NotifMention:    userSettings.NotifMention,
		NotifCommunity:  userSettings.NotifCommunity,
		NotifNewsletter: userSettings.NotifNewsletter,
	}

	returnData, err := anypb.New(userSettingsReturn)
	if err != nil {
		return &pb.ApiResponse{
			Success: false,
			Message: "An error occured: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponse{
		Success: true,
		Message: "Get User Settings successful.",
		Data:    returnData,
	}, nil
}

func (h *Handlers) User_UpdateSettings(ctx context.Context, req *pb.UpdateSettingsRequest) (*pb.ApiResponse, error) {
	zap.L().Info("User " + req.UserId + " is updating settings")

	var userSettings models.UserSetting
	err := h.DB.WithContext(ctx).Where("user_id = ?", req.UserId).First(&userSettings).Error
	if err != nil && err.Error() == "record not found" {
		userSettings = models.UserSetting{
			UserId:          req.UserId,
			FontSize:        "medium",
			FontColor:       "black",
			NotifLike:       true,
			NotifRepost:     true,
			NotifFollow:     true,
			NotifMention:    true,
			NotifCommunity:  true,
			NotifNewsletter: true,
		}

		err = h.DB.WithContext(ctx).Create(&userSettings).Error
		if err != nil {
			return &pb.ApiResponse{
				Success: false,
				Message: "An error occurred while creating user settings: " + err.Error(),
				Data:    nil,
			}, nil
		}
	} else {
		userSettings.FontSize = req.FontSize
		userSettings.FontColor = req.FontColor
		userSettings.NotifLike = req.NotifLike
		userSettings.NotifRepost = req.NotifRepost
		userSettings.NotifFollow = req.NotifFollow
		userSettings.NotifMention = req.NotifMention
		userSettings.NotifCommunity = req.NotifCommunity
		userSettings.NotifNewsletter = req.NotifNewsletter

		err = h.DB.WithContext(ctx).Save(&userSettings).Error
		if err != nil {
			return &pb.ApiResponse{
				Success: false,
				Message: "An error occurred while updating user settings: " + err.Error(),
				Data:    nil,
			}, nil
		}
	}

	return &pb.ApiResponse{
		Success: true,
		Message: "Update User Settings successful.",
		Data:    nil,
	}, nil
}
