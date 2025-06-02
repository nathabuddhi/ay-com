package handlers

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-user/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-user/proto/user"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"
)

func (h *Handlers) User_GetSettings(ctx context.Context, req *pb.GetSettingsRequest) (*pb.ApiResponseUser, error) {
	zap.L().Info("User " + req.UserId + " is getting settings")

	var userSettings models.UserSetting
	err := h.DB.WithContext(ctx).Where("user_id = ?", req.UserId).First(&userSettings).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		userSettings = models.UserSetting{
			UserId:    req.UserId,
			FontSize:  "medium",
			FontColor: "black",
		}

		err = h.DB.WithContext(ctx).Create(&userSettings).Error
		if err != nil {
			return &pb.ApiResponseUser{
				Success: false,
				Message: "An error occurred while creating user settings: " + err.Error(),
				Data:    nil,
			}, nil
		}

		var user models.User
		err = h.DB.WithContext(ctx).Where("user_id = ?", req.UserId).First(&user).Error
		if err != nil {
			return &pb.ApiResponseUser{
				Success: false,
				Message: "An error occurred while creating user settings: " + err.Error(),
				Data:    nil,
			}, nil
		}

		userSettingsReturn := &pb.UserSettings{
			FontSize:  "medium",
			FontColor: "black",
			Private:   false,
		}

		returnData, err := anypb.New(userSettingsReturn)
		if err != nil {
			return &pb.ApiResponseUser{
				Success: false,
				Message: "An error occured: " + err.Error(),
				Data:    nil,
			}, nil
		}

		return &pb.ApiResponseUser{
			Success: true,
			Message: "Get User Settings successful.",
			Data:    returnData,
		}, nil
	}

	var user models.User
	err = h.DB.WithContext(ctx).Where("user_id = ?", req.UserId).First(&user).Error
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "An error occurred while fetching user settings: " + err.Error(),
			Data:    nil,
		}, nil
	}
	userSettingsReturn := &pb.UserSettings{
		FontSize:  userSettings.FontSize,
		FontColor: userSettings.FontColor,
		Private:   user.IsPrivate,
	}

	returnData, err := anypb.New(userSettingsReturn)
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "An error occured: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Get User Settings successful.",
		Data:    returnData,
	}, nil
}

func (h *Handlers) User_UpdateSettings(ctx context.Context, req *pb.UpdateSettingsRequest) (*pb.ApiResponseUser, error) {
	zap.L().Info("User " + req.UserId + " is updating settings")

	var userSettings models.UserSetting
	err := h.DB.WithContext(ctx).Where("user_id = ?", req.UserId).First(&userSettings).Error
	if err != nil && err.Error() == "record not found" {
		userSettings = models.UserSetting{
			UserId:    req.UserId,
			FontSize:  "medium",
			FontColor: "black",
		}

		err = h.DB.WithContext(ctx).Create(&userSettings).Error
		if err != nil {
			return &pb.ApiResponseUser{
				Success: false,
				Message: "An error occurred while creating user settings: " + err.Error(),
				Data:    nil,
			}, nil
		}
	} else {
		userSettings.FontSize = req.FontSize
		userSettings.FontColor = req.FontColor

		err = h.DB.WithContext(ctx).Save(&userSettings).Error
		if err != nil {
			return &pb.ApiResponseUser{
				Success: false,
				Message: "An error occurred while updating user settings: " + err.Error(),
				Data:    nil,
			}, nil
		}
	}

	var user models.User
	err = h.DB.WithContext(ctx).Where("user_id = ?", req.UserId).First(&user).Error
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "An error occurred while updating user settings: " + err.Error(),
			Data:    nil,
		}, nil
	}
	user.IsPrivate = req.Private
	err = h.DB.WithContext(ctx).Save(&user).Error
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "An error occurred while updating user settings: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Update User Settings successful.",
		Data:    nil,
	}, nil
}
