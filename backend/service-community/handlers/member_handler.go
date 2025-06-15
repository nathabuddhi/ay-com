package handlers

import (
	"github.com/nathabuddhi/ay-com/backend/service-community/models"
)

func (h *Handler) IsUserModerator(userId string, communityId string) bool {
	var member models.CommunityMember
	err := h.DB.Where("user_id = ? AND community_id = ?", userId, communityId).First(&member).Error
	if err != nil {
		return false
	}

	return member.Role == "owner" || member.Role == "moderator"
}

func (h *Handler) IsUserMember(userId string, communityId string) bool {
	var member models.CommunityMember
	err := h.DB.Where("user_id = ? AND community_id = ?", userId, communityId).First(&member).Error
	if err != nil {
		return false
	}

	return member.Role == "owner" || member.Role == "moderator" || member.Role == "member"
}

func (h *Handler) IsUserOwner(userId string, communityId string) bool {
	var member models.CommunityMember
	err := h.DB.Where("user_id = ? AND community_id = ?", userId, communityId).First(&member).Error
	if err != nil {
		return false
	}

	return member.Role == "owner"
}

// func (h *Handler) Community_GetPendingUsers(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
// 	zap.L().Info("Getting all pending users for community", zap.String("community_id", req.Value))

// 	var pendingUsers []models.CommunityMember
// 	err := h.DB.WithContext(ctx).Where("community_id = ? AND status = ?", req.Value, "pending").Find(&pendingUsers).Error
// 	if err != nil {
// 		return &pb.ApiResponseCommunity{Success: false, Message: err.Error()}, nil
// 	}

// 	if len(pendingUsers) == 0 {
// 		return &pb.ApiResponseCommunity{Success: true, Message: "No pending users found."}, nil
// 	}

// 	pendingUserList := make([]*pb.CommunityMember, len(pendingUsers))
// 	for i, user := range pendingUsers {
// 		pendingUserList[i] = &pb.CommunityMember{
// 			UserId:   user.UserId,
// 			Role:     user.Role,
// 			JoinedAt: user.JoinedAt.String(),
// 		}
// 	}

// 	return &pb.ApiResponseCommunity{
// 		Success: true,
// 		Message: "Pending users retrieved successfully.",
// 		Data:    pendingUserList,
// 	}, nil
// }
