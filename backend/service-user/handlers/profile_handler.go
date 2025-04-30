package handlers

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-user/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-user/proto/user"
	"google.golang.org/protobuf/types/known/anypb"
)

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (h *Handlers) GetProfile(ctx context.Context, req *pb.GetUserRequest) (*pb.ApiResponse, error) {
	var user models.User
	if err := h.DB.Where("user_id = ?", req.UserId).First(&user).Error; err != nil {
		return &pb.ApiResponse{
			Success: false,
			Message: "User Not Found.",
		}, nil
	}

	if user.IsDeactivated || user.IsBanned {
		return &pb.ApiResponse{
			Success: false,
			Message: "User is inactive or currently banned.",
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
	}

	anyUser, err := anypb.New(userData)
	if err != nil {
		return &pb.ApiResponse{
			Success: false,
			Message: "An error occured: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponse{
		Success: true,
		Message: "Get User Profile successful.",
		Data:    anyUser,
	}, nil
}
