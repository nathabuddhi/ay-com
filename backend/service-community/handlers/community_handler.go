package handlers

import (
	"context"
	"time"

	"github.com/nathabuddhi/ay-com/backend/service-community/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-community/proto/community"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
)

func (h *Handler) Community_GetAllCommunities(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	zap.L().Info("Getting all communities")

	var communities []models.Community
	err := h.DB.WithContext(ctx).Where("is_pending = ?", false).Find(&communities).Error
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
		Message: "Get all communities successful.",
		Data:    returnData,
	}, nil
}

func (h *Handler) Community_CreateCommunity(ctx context.Context, req *pb.CreateCommunityRequest) (*pb.ApiResponseCommunity, error) {
	zap.L().Info("Creating community", zap.String("community_name", req.CommunityName))

	community := models.Community{
		CommunityId:   req.CommunityId,
		CommunityName: req.CommunityName,
		Description:   req.Description,
		Rules:         req.Rules,
		CreatorId:     req.UserId,
		IsPending:     true,
		CreatedAt:     time.Now(),
		BannerImage:   req.BannerImage,
		IconImage:     req.IconImage,
	}

	if err := h.DB.WithContext(ctx).Create(&community).Error; err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: err.Error()}, nil
	}

	owner := models.CommunityMember{
		CommunityId: req.CommunityId,
		UserId:      req.UserId,
		Role:        "owner",
		JoinedAt:    community.CreatedAt,
	}

	if err := h.DB.WithContext(ctx).Create(&owner).Error; err != nil {
		h.DB.WithContext(ctx).Delete(&community)
		return &pb.ApiResponseCommunity{Success: false, Message: "Failed to register community owner."}, nil
	}

	return &pb.ApiResponseCommunity{
		Success: true,
		Message: "Create community successful.",
		Data:    nil,
	}, nil
}
