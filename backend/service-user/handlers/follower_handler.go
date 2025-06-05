package handlers

import (
	"context"
	"encoding/json"

	"github.com/nathabuddhi/ay-com/backend/service-user/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-user/proto/user"
	"github.com/nathabuddhi/ay-com/backend/service-user/rabbitmq"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"
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

func (h *Handlers) User_FollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.ApiResponseUser, error) {
	zap.L().Info("User " + req.UserId + " is following user " + req.ToFollowId)

	if req.UserId == req.ToFollowId {
		return &pb.ApiResponseUser{Success: false, Message: "You cannot follow yourself."}, nil
	}

	var blockedUser models.BlockedUsers
	err := h.DB.WithContext(ctx).Where("user_id = ? AND blocked_id = ?", req.UserId, req.ToFollowId).First(&blockedUser).Error
	if err == nil {
		return &pb.ApiResponseUser{Success: false, Message: "You are blocked by this user and cannot follow them."}, nil
	}

	err = h.DB.WithContext(ctx).Where("user_id = ? AND blocked_id = ?", req.ToFollowId, req.UserId).First(&blockedUser).Error
	if err == nil {
		return &pb.ApiResponseUser{Success: false, Message: "This user has blocked you and you cannot follow them."}, nil
	}

	var userFollowing models.UserFollowing
	err = h.DB.WithContext(ctx).Where("user_id = ? AND followed_id = ?", req.UserId, req.ToFollowId).First(&userFollowing).Error
	if err == nil {
		return &pb.ApiResponseUser{Success: false, Message: "Already following this user."}, nil
	}
	userFollowing.UserId = req.UserId
	userFollowing.FollowedId = req.ToFollowId

	err = h.DB.WithContext(ctx).Create(&userFollowing).Error
	if err != nil {
		zap.L().Error("Failed to follow user: " + err.Error())
		return &pb.ApiResponseUser{Success: false, Message: "An unknown error occured. Please try again."}, nil
	}

	var user models.User
	var follower models.User
	h.DB.WithContext(ctx).Where("user_id = ?", req.ToFollowId).First(&user)
	h.DB.WithContext(ctx).Where("user_id = ?", req.UserId).First(&follower)
	rabbitmq.PublishSendNotification("follow", user.UserId, user.Email, "New Follower", follower.Username+" is now following you.", follower.UserId)
	rabbitmq.PublishDeleteRedis("getprofile/" + req.ToFollowId)
	rabbitmq.PublishDeleteRedis("getprofile/" + req.UserId)
	rabbitmq.PublishDeleteRedis("getallfollowers/" + req.ToFollowId)
	rabbitmq.PublishDeleteRedis("getallfollowing/" + req.UserId)
	return &pb.ApiResponseUser{Success: true, Message: "Followed user successfully."}, nil
}

func (h *Handlers) User_UnFollowUser(ctx context.Context, req *pb.UnFollowUserRequest) (*pb.ApiResponseUser, error) {
	zap.L().Info("User " + req.UserId + " is unfollowing user " + req.ToUnfollowId)

	deletedCount := h.DB.WithContext(ctx).Exec(`DELETE FROM user_followings WHERE user_id = ? AND followed_id = ?`, req.UserId, req.ToUnfollowId).RowsAffected

	if deletedCount == 0 {
		zap.L().Error("Failed to unfollow user: 0 rows affected.")
		return &pb.ApiResponseUser{Success: false, Message: "An unknown error occured. Please try again."}, nil
	}

	rabbitmq.PublishDeleteRedis("getprofile/" + req.ToUnfollowId)
	rabbitmq.PublishDeleteRedis("getprofile/" + req.UserId)
	return &pb.ApiResponseUser{Success: true, Message: "Unfollowed user successfully."}, nil
}

func (h *Handlers) User_BlockUser(ctx context.Context, req *pb.BlockUserRequest) (*pb.ApiResponseUser, error) {
	zap.L().Info("User " + req.UserId + " is blocking user " + req.ToBlockId)

	if req.UserId == req.ToBlockId {
		return &pb.ApiResponseUser{Success: false, Message: "You cannot block yourself."}, nil
	}

	var userFollower models.UserFollowing
	unFollowedCount := h.DB.WithContext(ctx).Where("user_id = ? AND followed_id = ?", req.ToBlockId, req.UserId).Delete(&userFollower).RowsAffected
	unFollowedCount += h.DB.WithContext(ctx).Where("user_id = ? AND followed_id = ?", req.UserId, req.ToBlockId).Delete(&userFollower).RowsAffected

	if unFollowedCount != 0 {
		rabbitmq.PublishDeleteRedis("getprofile/" + req.ToBlockId)
		rabbitmq.PublishDeleteRedis("getprofile/" + req.UserId)
		rabbitmq.PublishDeleteRedis("getallfollowers/" + req.ToBlockId)
		rabbitmq.PublishDeleteRedis("getallfollowers/" + req.UserId)
		rabbitmq.PublishDeleteRedis("getallfollowing/" + req.ToBlockId)
		rabbitmq.PublishDeleteRedis("getallfollowing/" + req.UserId)
	}

	var blockedUser models.BlockedUsers
	err := h.DB.WithContext(ctx).Where("user_id = ? AND blocked_id = ?", req.UserId, req.ToBlockId).First(&blockedUser).Error
	if err == nil {
		return &pb.ApiResponseUser{Success: false, Message: "Already blocking this user."}, nil
	}
	blockedUser.UserId = req.UserId
	blockedUser.BlockedId = req.ToBlockId

	err = h.DB.WithContext(ctx).Create(&blockedUser).Error
	if err != nil {
		zap.L().Error("Failed to block user: " + err.Error())
		return &pb.ApiResponseUser{Success: false, Message: "An unknown error occured. Please try again."}, nil
	}
	return &pb.ApiResponseUser{Success: true, Message: "Blocked user successfully."}, nil
}

func (h *Handlers) User_UnBlockUser(ctx context.Context, req *pb.UnBlockUserRequest) (*pb.ApiResponseUser, error) {
	zap.L().Info("User " + req.UserId + " is unblocking user " + req.ToUnblockId)

	var userFollower models.BlockedUsers
	deletedCount := h.DB.WithContext(ctx).Where("user_id = ? AND blocked_id = ?", req.UserId, req.ToUnblockId).Delete(&userFollower).RowsAffected

	if deletedCount == 0 {
		zap.L().Error("Failed to unblock user: 0 rows affected.")
		return &pb.ApiResponseUser{Success: false, Message: "An unknown error occured. Please try again."}, nil
	}

	return &pb.ApiResponseUser{Success: true, Message: "Unblocked user successfully."}, nil
}

func (h *Handlers) User_GetAllFollowers(ctx context.Context, req *pb.GetAllFollowersRequest) (*pb.ApiResponseUser, error) {
	zap.L().Info("User " + req.RequesterId + " is getting all followers of user " + req.UserId)

	var user models.User
	err := h.DB.WithContext(ctx).Where("user_id = ?", req.UserId).First(&user).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "User not found."}, nil
	}

	var requester models.User
	err = h.DB.WithContext(ctx).Where("user_id = ?", req.RequesterId).First(&requester).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Requester not found."}, nil
	}

	if requester.IsBanned || requester.IsDeactivated {
		return &pb.ApiResponseUser{Success: false, Message: "Your account is not active or is banned."}, nil
	}

	if user.IsBanned || user.IsDeactivated {
		return &pb.ApiResponseUser{Success: false, Message: "This user account is not active or is banned."}, nil
	}

	if user.IsPrivate && req.UserId != req.RequesterId {
		var userFollowing models.UserFollowing
		err = h.DB.WithContext(ctx).Where("user_id = ? AND followed_id = ?", req.RequesterId, req.UserId).First(&userFollowing).Error
		if err != nil {
			return &pb.ApiResponseUser{Success: false, Message: "This user account is private."}, nil
		}
	}

	var followers []models.UserFollowing
	err = h.DB.WithContext(ctx).Where("followed_id = ?", req.UserId).Find(&followers).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "An unknown error occured. Please try again."}, nil
	}

	allFollowersResponse := &pb.AllFollowersResponse{
		Followers: make([]*pb.StringUser, len(followers)),
	}

	for i, request := range followers {
		allFollowersResponse.Followers[i] = &pb.StringUser{
			Value: request.UserId,
		}
	}

	if !user.IsPrivate {
		redisData, err := json.Marshal(allFollowersResponse)
		if err == nil {
			rabbitmq.PublishSetRedis("getallfollowers/"+user.UserId, string(redisData))
		}
	}

	returnData, err := anypb.New(allFollowersResponse)
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "An error occured: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Get All Followers successful.",
		Data:    returnData,
	}, nil
}

func (h *Handlers) User_GetAllFollowing(ctx context.Context, req *pb.GetAllFollowingRequest) (*pb.ApiResponseUser, error) {
	zap.L().Info("User " + req.RequesterId + " is getting all followings of user " + req.UserId)

	var user models.User
	err := h.DB.WithContext(ctx).Where("user_id = ?", req.UserId).First(&user).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "User not found."}, nil
	}

	var requester models.User
	err = h.DB.WithContext(ctx).Where("user_id = ?", req.RequesterId).First(&requester).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Requester not found."}, nil
	}

	if requester.IsBanned || requester.IsDeactivated {
		return &pb.ApiResponseUser{Success: false, Message: "Your account is not active or is banned."}, nil
	}

	if user.IsBanned || user.IsDeactivated {
		return &pb.ApiResponseUser{Success: false, Message: "This user account is not active or is banned."}, nil
	}

	if user.IsPrivate && req.UserId != req.RequesterId {
		var userFollowing models.UserFollowing
		err = h.DB.WithContext(ctx).Where("user_id = ? AND followed_id = ?", req.RequesterId, req.UserId).First(&userFollowing).Error
		if err != nil {
			return &pb.ApiResponseUser{Success: false, Message: "This user account is private."}, nil
		}
	}

	var followers []models.UserFollowing
	err = h.DB.WithContext(ctx).Where("user_id = ?", req.UserId).Find(&followers).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "An unknown error occured. Please try again."}, nil
	}

	allFollowingResponse := &pb.AllFollowingResponse{
		Following: make([]*pb.StringUser, len(followers)),
	}

	for i, request := range followers {
		allFollowingResponse.Following[i] = &pb.StringUser{
			Value: request.FollowedId,
		}
	}

	if !user.IsPrivate {
		redisData, err := json.Marshal(allFollowingResponse)
		if err == nil {
			rabbitmq.PublishSetRedis("getallfollowing/"+user.UserId, string(redisData))
		}
	}

	returnData, err := anypb.New(allFollowingResponse)
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "An error occured: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Get All Followings successful.",
		Data:    returnData,
	}, nil
}

func (h *Handlers) User_GetAllBlocked(ctx context.Context, req *pb.StringUser) (*pb.ApiResponseUser, error) {
	zap.L().Info("User " + req.Value + " is getting all blocked users.")

	var user models.User
	err := h.DB.WithContext(ctx).Where("user_id = ?", req.Value).First(&user).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "User not found."}, nil
	}

	if user.IsBanned || user.IsDeactivated {
		return &pb.ApiResponseUser{Success: false, Message: "This user account is not active or is banned."}, nil
	}

	var blocked []models.BlockedUsers
	err = h.DB.WithContext(ctx).Where("user_id = ?", req.Value).Find(&blocked).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "An unknown error occured. Please try again."}, nil
	}

	allBlockedResponse := &pb.AllBlockedUserResponse{
		Blocked: make([]*pb.BlockedUser, len(blocked)),
	}

	for i, request := range blocked {
		var blockedUser models.User
		err = h.DB.WithContext(ctx).Where("user_id = ?", request.BlockedId).First(&blockedUser).Error

		if err != nil {
			zap.L().Error("Failed to get blocked user: " + err.Error())
		} else {
			zap.L().Info("Blocked User: ", zap.Any("blocked", request))
			allBlockedResponse.Blocked[i] = &pb.BlockedUser{
				UserId:   blockedUser.UserId,
				Username: blockedUser.Username,
				Name:     blockedUser.Name,
			}
		}
	}
	returnData, err := anypb.New(allBlockedResponse)
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "An error occured: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Get All Blocked successful.",
		Data:    returnData,
	}, nil
}

func (h *Handlers) User_IsUserFollowing(ctx context.Context, req *pb.IsUserFollowingRequest) (*pb.BoolUser, error) {
	zap.L().Info("User " + req.FollowerId + " is checking if they are following user " + req.FollowingId)

	var userFollowing models.UserFollowing
	err := h.DB.WithContext(ctx).Where("user_id = ? AND followed_id = ?", req.FollowerId, req.FollowingId).First(&userFollowing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.BoolUser{Value: false}, nil
		}
		return &pb.BoolUser{Value: false}, nil
	}

	return &pb.BoolUser{Value: true}, nil
}
