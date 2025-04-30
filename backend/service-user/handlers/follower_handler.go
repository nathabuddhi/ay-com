package handlers

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-user/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-user/proto/user"
	"github.com/nathabuddhi/ay-com/backend/service-user/rabbitmq"
	"go.uber.org/zap"
)

func (h *Handlers) GetFollowing(user_id string) (int, error) {
	var count int64
	err := h.DB.Model(&models.UserFollowing{}).Where("user_id = ?", user_id).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (h *Handlers) GetFollowers(user_id string) (int, error) {
	var count int64
	err := h.DB.Model(&models.UserFollowing{}).Where("followed_id = ?", user_id).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (h *Handlers) User_FollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.ApiResponse, error) {
	zap.L().Info("User " + req.UserId + " is following user " + req.ToFollowId)

	if req.UserId == req.ToFollowId {
		return &pb.ApiResponse{Success: false, Message: "You cannot follow yourself."}, nil
	}

	var blockedUser models.BlockedUsers
	err := h.DB.WithContext(ctx).Where("user_id = ? AND blocked_id = ?", req.UserId, req.ToFollowId).First(&blockedUser).Error
	if err == nil {
		return &pb.ApiResponse{Success: false, Message: "You are blocked by this user and cannot follow them."}, nil
	}

	err = h.DB.WithContext(ctx).Where("user_id = ? AND blocked_id = ?", req.ToFollowId, req.UserId).First(&blockedUser).Error
	if err == nil {
		return &pb.ApiResponse{Success: false, Message: "This user has blocked you and you cannot follow them."}, nil
	}

	var userFollowing models.UserFollowing
	err = h.DB.WithContext(ctx).Where("user_id = ? AND followed_id = ?", req.UserId, req.ToFollowId).First(&userFollowing).Error
	if err == nil {
		return &pb.ApiResponse{Success: false, Message: "Already following this user."}, nil
	}
	userFollowing.UserId = req.UserId
	userFollowing.FollowedId = req.ToFollowId

	err = h.DB.WithContext(ctx).Create(&userFollowing).Error
	if err != nil {
		zap.L().Error("Failed to follow user: " + err.Error())
		return &pb.ApiResponse{Success: false, Message: "An unknown error occured. Please try again."}, nil
	}
	rabbitmq.PublishDeleteRedis("getprofile/" + req.ToFollowId)
	rabbitmq.PublishDeleteRedis("getprofile/" + req.UserId)
	return &pb.ApiResponse{Success: true, Message: "Followed user successfully."}, nil
}

func (h *Handlers) User_UnFollowUser(ctx context.Context, req *pb.UnFollowUserRequest) (*pb.ApiResponse, error) {
	zap.L().Info("User " + req.UserId + " is unfollowing user " + req.ToUnfollowId)

	deletedCount := h.DB.WithContext(ctx).Exec(`DELETE FROM user_followings WHERE user_id = ? AND followed_id = ?`, req.UserId, req.ToUnfollowId).RowsAffected

	if deletedCount == 0 {
		zap.L().Error("Failed to unfollow user: 0 rows affected.")
		return &pb.ApiResponse{Success: false, Message: "An unknown error occured. Please try again."}, nil
	}

	rabbitmq.PublishDeleteRedis("getprofile/" + req.ToUnfollowId)
	rabbitmq.PublishDeleteRedis("getprofile/" + req.UserId)
	return &pb.ApiResponse{Success: true, Message: "Unfollowed user successfully."}, nil
}

func (h *Handlers) User_BlockUser(ctx context.Context, req *pb.BlockUserRequest) (*pb.ApiResponse, error) {
	zap.L().Info("User " + req.UserId + " is blocking user " + req.ToBlockId)

	if req.UserId == req.ToBlockId {
		return &pb.ApiResponse{Success: false, Message: "You cannot block yourself."}, nil
	}

	var userFollower models.UserFollowing
	unFollowedCount := h.DB.WithContext(ctx).Where("user_id = ? AND followed_id = ?", req.ToBlockId, req.UserId).Delete(&userFollower).RowsAffected
	unFollowedCount += h.DB.WithContext(ctx).Where("user_id = ? AND followed_id = ?", req.UserId, req.ToBlockId).Delete(&userFollower).RowsAffected

	if unFollowedCount != 0 {
		rabbitmq.PublishDeleteRedis("getprofile/" + req.ToBlockId)
		rabbitmq.PublishDeleteRedis("getprofile/" + req.UserId)
	}

	var blockedUser models.BlockedUsers
	err := h.DB.WithContext(ctx).Where("user_id = ? AND blocked_id = ?", req.UserId, req.ToBlockId).First(&blockedUser).Error
	if err == nil {
		return &pb.ApiResponse{Success: false, Message: "Already blocking this user."}, nil
	}
	blockedUser.UserId = req.UserId
	blockedUser.BlockedId = req.ToBlockId

	err = h.DB.WithContext(ctx).Create(&blockedUser).Error
	if err != nil {
		zap.L().Error("Failed to block user: " + err.Error())
		return &pb.ApiResponse{Success: false, Message: "An unknown error occured. Please try again."}, nil
	}
	return &pb.ApiResponse{Success: true, Message: "Blocked user successfully."}, nil
}

func (h *Handlers) User_UnBlockUser(ctx context.Context, req *pb.UnBlockUserRequest) (*pb.ApiResponse, error) {
	zap.L().Info("User " + req.UserId + " is unblocking user " + req.ToUnblockId)

	var userFollower models.BlockedUsers
	deletedCount := h.DB.WithContext(ctx).Where("user_id = ? AND blocked_id = ?", req.UserId, req.ToUnblockId).Delete(&userFollower).RowsAffected

	if deletedCount == 0 {
		zap.L().Error("Failed to unblock user: 0 rows affected.")
		return &pb.ApiResponse{Success: false, Message: "An unknown error occured. Please try again."}, nil
	}

	return &pb.ApiResponse{Success: true, Message: "Unblocked user successfully."}, nil
}
