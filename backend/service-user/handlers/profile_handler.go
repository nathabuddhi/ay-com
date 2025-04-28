package handlers

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-user/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-user/proto/user"
	"google.golang.org/protobuf/types/known/anypb"
)

func (h *Handlers) GetProfile(ctx context.Context, req *pb.GetUserRequest) (*pb.ApiResponse, error) {
	var user models.User
	if err := h.DB.First(&user, "id = ?", req.UserId).Error; err != nil {
		return &pb.ApiResponse{
			Success: false,
			Message: "User not found",
		}, err
	}

	userData := &pb.UserProfile{
		UserId:         user.UserId,
		Username:       user.Username,
		IsVerified:     user.IsVerified,
		Name:           user.Name,
		Banner:         *user.Banner,
		ProfilePicture: *user.ProfilePicture,
		Bio:            *user.Bio,
		Followers:      0,
		Following:      0,
	}

	anyUser, err := anypb.New(userData)
	if err != nil {
		return nil, err
	}

	return &pb.ApiResponse{
		Success: true,
		Message: "Get User Profile successful.",
		Data:    anyUser,
	}, nil
}
