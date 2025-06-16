package handlers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nathabuddhi/ay-com/backend/service-community/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-community/proto/community"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
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

func (h *Handler) GetUserRole(userId string, communityId string) string {
	var member models.CommunityMember
	err := h.DB.Where("user_id = ? AND community_id = ?", userId, communityId).First(&member).Error
	if err != nil {
		return "guest"
	}

	return member.Role
}

func (h *Handler) IsUserPendingJoin(userId string, communityId string) bool {
	var joinRequest models.CommunityJoinRequest
	err := h.DB.Where("user_id = ? AND community_id = ?", userId, communityId).First(&joinRequest).Error
	if err != nil {
		return false
	}
	if joinRequest.Status == "pending" {
		return true
	}
	return false
}

func (h *Handler) Community_GetCommunityMembers(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	var members []models.CommunityMember
	err := h.DB.WithContext(ctx).
		Where("community_id = ?", req.Value).
		Order("role DESC, joined_at DESC").
		Find(&members).Error
	if err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: "Failed to get community members: " + err.Error()}, nil
	}

	memberReturn := &pb.GetMembersResponse{
		Members: make([]*pb.CommunityMember, len(members)),
	}

	for i, member := range members {
		memberReturn.Members[i] = &pb.CommunityMember{
			UserId: member.UserId,
			Role:   member.Role,
		}
	}

	returnData, err := anypb.New(memberReturn)

	if err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: "Failed to marshal response data: " + err.Error()}, nil
	}

	return &pb.ApiResponseCommunity{
		Success: true,
		Message: "Community members retrieved successfully.",
		Data:    returnData,
	}, nil
}

func (h *Handler) Community_GetPendingUsers(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	zap.L().Info("Getting all pending users for community", zap.String("community_id", req.Value))

	var pendingUsers []models.CommunityJoinRequest
	err := h.DB.WithContext(ctx).Where("community_id = ? AND status = ?", req.Value, "pending").Find(&pendingUsers).Error
	if err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: err.Error()}, nil
	}

	if len(pendingUsers) == 0 {
		return &pb.ApiResponseCommunity{Success: true, Message: "No pending users found."}, nil
	}

	pendingUserList := make([]*pb.CommunityMember, len(pendingUsers))
	for i, user := range pendingUsers {
		pendingUserList[i] = &pb.CommunityMember{
			UserId: user.UserId,
			Role:   "pending",
		}
	}

	returnData, err := anypb.New(&pb.GetMembersResponse{
		Members: pendingUserList,
	})

	if err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: "Failed to marshal response data: " + err.Error()}, nil
	}

	return &pb.ApiResponseCommunity{
		Success: true,
		Message: "Pending users retrieved successfully.",
		Data:    returnData,
	}, nil
}

func (h *Handler) Community_JoinCommunity(ctx context.Context, req *pb.GeneralCommunityRequest) (*pb.ApiResponseCommunity, error) {
	if req.CommunityId == "" || req.UserId == "" {
		return &pb.ApiResponseCommunity{Success: false, Message: "Community ID and User ID cannot be empty."}, nil
	}

	if h.IsUserMember(req.UserId, req.CommunityId) {
		return &pb.ApiResponseCommunity{Success: false, Message: "User is already a member of the community."}, nil

	}
	var joinRequest models.CommunityJoinRequest
	err := h.DB.WithContext(ctx).Where("community_id = ? AND user_id = ?", req.CommunityId, req.UserId).First(&joinRequest).Error
	if err == nil {
		if joinRequest.Status == "pending" {
			return &pb.ApiResponseCommunity{Success: false, Message: "User has already requested to join the community."}, nil
		}
	}
	joinRequest = models.CommunityJoinRequest{
		RequestId:   uuid.New().String(),
		CommunityId: req.CommunityId,
		UserId:      req.UserId,
		Status:      "pending",
		RequestedAt: time.Now(),
	}
	if err := h.DB.WithContext(ctx).Create(&joinRequest).Error; err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: "Failed to create join request: " + err.Error()}, nil
	}
	return &pb.ApiResponseCommunity{
		Success: true,
		Message: "Join request created successfully. Awaiting approval.",
		Data:    nil,
	}, nil
}

func (h *Handler) Community_ApproveJoinRequest(ctx context.Context, req *pb.GeneralModeratorRequest) (*pb.ApiResponseCommunity, error) {
	if req.CommunityId == "" || req.UserId == "" {
		return &pb.ApiResponseCommunity{Success: false, Message: "Community ID and User ID cannot be empty."}, nil
	}

	if !h.IsUserModerator(req.ModeratorId, req.CommunityId) {
		return &pb.ApiResponseCommunity{Success: false, Message: "Only moderators can approve join requests."}, nil
	}

	var joinRequest models.CommunityJoinRequest
	err := h.DB.WithContext(ctx).Where("community_id = ? AND user_id = ?", req.CommunityId, req.UserId).First(&joinRequest).Error
	if err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: "Join request not found: " + err.Error()}, nil
	}

	if joinRequest.Status != "pending" {
		return &pb.ApiResponseCommunity{Success: false, Message: "Join request is not pending."}, nil
	}

	joinRequest.Status = "accepted"
	if err := h.DB.WithContext(ctx).Save(&joinRequest).Error; err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: "Failed to approve join request: " + err.Error()}, nil
	}

	member := models.CommunityMember{
		CommunityId: req.CommunityId,
		UserId:      req.UserId,
		Role:        "member",
	}
	if err := h.DB.WithContext(ctx).Create(&member).Error; err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: "Failed to add member to community: " + err.Error()}, nil
	}

	return &pb.ApiResponseCommunity{
		Success: true,
		Message: "Join request approved successfully.",
		Data:    nil,
	}, nil
}

func (h *Handler) Community_DenyJoinRequest(ctx context.Context, req *pb.GeneralModeratorRequest) (*pb.ApiResponseCommunity, error) {
	if req.CommunityId == "" || req.UserId == "" {
		return &pb.ApiResponseCommunity{Success: false, Message: "Community ID and User ID cannot be empty."}, nil
	}

	if !h.IsUserModerator(req.ModeratorId, req.CommunityId) {
		return &pb.ApiResponseCommunity{Success: false, Message: "Only moderators can reject join requests."}, nil
	}

	var joinRequest models.CommunityJoinRequest
	err := h.DB.WithContext(ctx).Where("community_id = ? AND user_id = ?", req.CommunityId, req.UserId).First(&joinRequest).Error
	if err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: "Join request not found: " + err.Error()}, nil
	}
	if joinRequest.Status != "pending" {
		return &pb.ApiResponseCommunity{Success: false, Message: "Join request is not pending."}, nil
	}
	joinRequest.Status = "rejected"
	if err := h.DB.WithContext(ctx).Save(&joinRequest).Error; err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: "Failed to reject join request: " + err.Error()}, nil
	}
	return &pb.ApiResponseCommunity{
		Success: true,
		Message: "Join request rejected successfully.",
		Data:    nil,
	}, nil
}

func (h *Handler) Community_DemoteMember(ctx context.Context, req *pb.GeneralModeratorRequest) (*pb.ApiResponseCommunity, error) {
	if req.CommunityId == "" || req.UserId == "" {
		return &pb.ApiResponseCommunity{Success: false, Message: "Community ID and User ID cannot be empty."}, nil
	}

	if !h.IsUserOwner(req.ModeratorId, req.CommunityId) {
		return &pb.ApiResponseCommunity{Success: false, Message: "Only owners can demote moderators."}, nil
	}

	var member models.CommunityMember
	err := h.DB.WithContext(ctx).Where("community_id = ? AND user_id = ? AND role = ?", req.CommunityId, req.UserId, "moderator").First(&member).Error
	if err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: "Member not found: " + err.Error()}, nil
	}

	if member.Role == "owner" {
		return &pb.ApiResponseCommunity{Success: false, Message: "Cannot demote the owner of the community."}, nil
	}

	member.Role = "member"
	if err := h.DB.WithContext(ctx).Save(&member).Error; err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: "Failed to demote member: " + err.Error()}, nil
	}

	return &pb.ApiResponseCommunity{
		Success: true,
		Message: "Member demoted successfully.",
		Data:    nil,
	}, nil
}

func (h *Handler) Community_PromoteMember(ctx context.Context, req *pb.GeneralModeratorRequest) (*pb.ApiResponseCommunity, error) {
	if req.CommunityId == "" || req.UserId == "" {
		return &pb.ApiResponseCommunity{Success: false, Message: "Community ID and User ID cannot be empty."}, nil
	}

	if !h.IsUserOwner(req.ModeratorId, req.CommunityId) {
		return &pb.ApiResponseCommunity{Success: false, Message: "Only owners can promote members."}, nil
	}

	var member models.CommunityMember
	err := h.DB.WithContext(ctx).Where("community_id = ? AND user_id = ? AND role = ?", req.CommunityId, req.UserId, "member").First(&member).Error
	if err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: "Member not found: " + err.Error()}, nil
	}

	if member.Role == "owner" {
		return &pb.ApiResponseCommunity{Success: false, Message: "Cannot promote the owner of the community."}, nil
	}

	member.Role = "moderator"
	if err := h.DB.WithContext(ctx).Save(&member).Error; err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: "Failed to promote member: " + err.Error()}, nil
	}

	return &pb.ApiResponseCommunity{
		Success: true,
		Message: "Member promoted successfully.",
		Data:    nil,
	}, nil
}
