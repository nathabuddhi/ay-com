package handlers

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-community/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-community/proto/community"
	"github.com/nathabuddhi/ay-com/backend/service-community/rabbitmq"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
)

func (h *Handler) Admin_GetAllCommunityRequests(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	zap.L().Info("Admin getting all community requests")

	var communities []models.Community
	err := h.DB.WithContext(ctx).Where("is_pending = ?", true).Find(&communities).Error
	if err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: "Failed to get communities."}, nil
	}

	if len(communities) == 0 {
		return &pb.ApiResponseCommunity{Success: true, Message: "No communities found."}, nil
	}

	getAllCommunitiesResponse := &pb.GetCommunitiesResponse{
		Communities: make([]*pb.Community, len(communities)),
	}

	for i, community := range communities {
		getAllCommunitiesResponse.Communities[i], _ = h.GetCommunityById(community.CommunityId)
	}

	returnData, err := anypb.New(getAllCommunitiesResponse)
	if err != nil {
		return &pb.ApiResponseCommunity{
			Success: false,
			Message: "An error occurred: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponseCommunity{
		Success: true,
		Message: "Get all community requests successful.",
		Data:    returnData,
	}, nil
}

func (h *Handler) Admin_ApproveCommunity(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	zap.L().Info("Admin approving community", zap.String("community_id", req.Value))

	var community models.Community
	err := h.DB.WithContext(ctx).Where("community_id = ? AND is_pending = ?", req.Value, true).First(&community).Error
	if err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: "Community not found or already approved."}, nil
	}
	community.IsPending = false
	err = h.DB.WithContext(ctx).Save(&community).Error
	if err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: "Failed to approve community."}, nil
	}

	rabbitmq.PublishSendNotification("system", community.CreatorId, "", "Your new community "+community.CommunityName+" has been approved.", "Your previous community request has been approved and has been created. Your community is live and can now be accessed by all users within AY.com. Thankyou.", "system")

	return &pb.ApiResponseCommunity{
		Success: true,
		Message: "Community approved successfully.",
	}, nil
}

func (h *Handler) Admin_RejectCommunity(ctx context.Context, req *pb.RejectCommunityRequest) (*pb.ApiResponseCommunity, error) {
	zap.L().Info("Admin rejecting community", zap.String("community_id", req.CommunityId))

	var community models.Community
	err := h.DB.WithContext(ctx).Where("community_id = ? AND is_pending = ?", req.CommunityId, true).First(&community).Error
	if err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: "Community not found or already approved."}, nil
	}

	community.IsPending = false
	community.IsRejected = true
	community.RejectReason = req.Reason
	err = h.DB.WithContext(ctx).Save(&community).Error
	if err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: "Failed to reject community."}, nil
	}

	rabbitmq.PublishSendNotification("system", community.CreatorId, "", "Your community request "+community.CommunityName+" has been rejected.", "Your previous community request has been rejected due to "+req.Reason+"and has been deleted. You may retry at any time, thankyou.", "system")

	return &pb.ApiResponseCommunity{
		Success: true,
		Message: "Community approved successfully.",
	}, nil
}
