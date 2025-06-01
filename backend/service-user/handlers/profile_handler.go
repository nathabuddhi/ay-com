package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nathabuddhi/ay-com/backend/service-user/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-user/proto/user"
	"github.com/nathabuddhi/ay-com/backend/service-user/rabbitmq"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/anypb"
)

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (h *Handlers) User_GetProfile(ctx context.Context, req *pb.GetProfileRequest, isUserId bool) (*pb.ApiResponseUser, error) {
	var user models.User

	if isUserId {
		if err := h.DB.Where("user_id = ?", req.UserId).First(&user).Error; err != nil {
			return &pb.ApiResponseUser{
				Success: false,
				Message: "User Not Found.",
			}, nil
		}
	} else {
		if err := h.DB.Where("username = ?", req.UserId).First(&user).Error; err != nil {
			return &pb.ApiResponseUser{
				Success: false,
				Message: "User Not Found.",
			}, nil
		}
	}

	var requester models.User
	if err := h.DB.Where("user_id = ?", req.RequesterId).First(&requester).Error; err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "Requester Not Found.",
		}, nil
	}

	if user.IsDeactivated || user.IsBanned {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "User is inactive or currently banned.",
		}, nil
	}

	if requester.IsDeactivated || requester.IsBanned {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "Your account is inactive or currently banned.",
		}, nil
	}

	followers, err := h.GetFollowers(user.UserId)
	following, err2 := h.GetFollowing(user.UserId)

	if err != nil {
		followers = 0
	} else if err2 != nil {
		following = 0
	}

	bio := safeString(user.Bio)

	userData := &pb.UserProfile{
		UserId:     user.UserId,
		Username:   user.Username,
		IsVerified: user.IsVerified,
		Name:       user.Name,
		Bio:        bio,
		Followers:  int32(followers),
		Following:  int32(following),
		JoinDate:   user.JoinedAt.Format("2006-01-02"),
	}

	zap.L().Info(user.JoinedAt.Format("2006-01-02"))

	returnData, err := anypb.New(userData)
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "An error occured: " + err.Error(),
			Data:    nil,
		}, nil
	}

	redisData, err := json.Marshal(userData)
	if err == nil {
		rabbitmq.PublishSetRedis("getprofile/"+user.UserId, string(redisData))
	}

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Get User Profile successful.",
		Data:    returnData,
	}, nil
}

func (h *Handlers) User_DeactivateAccount(ctx context.Context, req *pb.DeactivateAccountRequest) (*pb.ApiResponseUser, error) {
	zap.L().Info("User Deactivating Account", zap.String("user_id", req.UserId))

	var user models.User
	if err := h.DB.Where("user_id = ?", req.UserId).First(&user).Error; err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "User Not Found.",
		}, nil
	}

	if user.IsDeactivated {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "User is already deactivated.",
		}, nil
	}

	if user.IsBanned {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "User is banned.",
		}, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "Invalid Credentials.",
		}, nil
	}

	err := h.DB.Model(&models.User{}).
		Where("user_id = ?", req.UserId).
		Update("is_deactivated", true).Error
	if err != nil {
		zap.L().Error("Failed to deactivate user account", zap.Error(err))
		return &pb.ApiResponseUser{
			Success: false,
			Message: "Failed to deactivate user account:" + err.Error(),
		}, nil
	}

	rabbitmq.PublishEmail(user.Email,
		"AY.com Account Deactivation",
		"Your account has been deactivated. If this was a mistake, please head to the activation page or contact support.",
	)

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Account deactivated successfully.",
	}, nil
}

func (h *Handlers) User_UpdateProfile(ctx context.Context, req *pb.UpdateUserProfileRequest) (*pb.ApiResponseUser, error) {
	zap.L().Info("User Updating Profile", zap.String("user_id", req.UserId))

	if req.Username == "" || req.Name == "" || req.DateOfBirth == "" || req.Bio == "" || req.Gender == "" {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "All fields are required.",
		}, nil
	}

	var user models.User
	if err := h.DB.Where("user_id = ?", req.UserId).First(&user).Error; err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "User Not Found.",
		}, nil
	}

	if user.IsDeactivated {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "User is already deactivated.",
		}, nil
	}

	if user.IsBanned {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "User is banned.",
		}, nil
	}

	user.Username = req.Username

	if err := h.DB.Where("username = ?", req.Username).First(&user).Error; err != nil && user.Username != req.Username {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "Username already taken.",
		}, nil
	}

	user.Name = req.Name
	user.Bio = &req.Bio
	parsedDate, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "Invalid date format. Please use YYYY-MM-DD.",
		}, nil
	}

	if time.Since(parsedDate) < 13*365*24*time.Hour {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "User must be at least 13 years old.",
		}, nil
	}
	user.DateOfBirth = parsedDate
	user.Gender = req.Gender

	if err := h.DB.Save(&user).Error; err != nil {
		zap.L().Error("Failed to update user profile", zap.Error(err))
		return &pb.ApiResponseUser{
			Success: false,
			Message: "Failed to update user profile: " + err.Error(),
		}, nil
	}

	return &pb.ApiResponseUser{
		Success: true,
		Message: "User profile updated successfully.",
	}, nil
}

func (h *Handlers) User_GetSelfProfile(ctx context.Context, req *pb.StringUser) (*pb.ApiResponseUser, error) {
	var user models.User
	if err := h.DB.Where("user_id = ?", req.Value).First(&user).Error; err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "User Not Found.",
		}, nil
	}

	if user.IsDeactivated || user.IsBanned {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "User is inactive or currently banned.",
		}, nil
	}

	bio := safeString(user.Bio)

	userData := &pb.UserProfile{
		UserId:      user.UserId,
		Username:    user.Username,
		Name:        user.Name,
		Bio:         bio,
		Gender:      user.Gender,
		DateOfBirth: user.DateOfBirth.Format("2006-01-02"),
		Email:       user.Email,
		JoinDate:    user.JoinedAt.Format("2006-01-02"),
		IsVerified:  user.IsVerified,
	}

	returnData, err := anypb.New(userData)
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "An error occured: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Get Self Profile successful.",
		Data:    returnData,
	}, nil
}
