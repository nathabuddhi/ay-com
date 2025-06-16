package handlers

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-community/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-community/proto/community"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
)

func (h *Handler) Community_GetCategories(ctx context.Context, req *pb.StringCommunity) (*pb.GetCategoriesResponse, error) {
	zap.L().Info("Getting all categories")

	var categories []models.CommunityCategory
	err := h.DB.WithContext(ctx).Find(&categories).Error
	if err != nil {
		return &pb.GetCategoriesResponse{Categories: []string{"General"}}, nil
	}

	if len(categories) == 0 {
		return &pb.GetCategoriesResponse{Categories: []string{"General"}}, nil
	}
	getCategoriesResponse := &pb.GetCategoriesResponse{
		Categories: make([]string, len(categories)),
	}

	for i, category := range categories {
		getCategoriesResponse.Categories[i] = category.Category
	}

	return getCategoriesResponse, nil
}

func (h *Handler) GetCommunityCategories(communityId string) []string {
	var categories []models.CommunityCategoryRelation
	err := h.DB.Where("community_id = ?", communityId).Find(&categories).Error
	if err != nil {
		zap.L().Error("Error fetching community categories", zap.Error(err))
		return []string{"General"}
	}

	categoryList := make([]string, len(categories))
	for i, category := range categories {
		categoryList[i] = category.Category
	}

	return categoryList
}

func (h *Handler) Community_GetAllUserPendingApprovalCommunities(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	zap.L().Info("Getting all user pending approval communities", zap.String("user_id", req.Value))

	var communities []models.Community
	err := h.DB.WithContext(ctx).Where("user_id = ? AND is_pending = ? AND is_rejected = ?", req.Value, true, false).Find(&communities).Error
	if err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: err.Error()}, nil
	}

	if len(communities) == 0 {
		return &pb.ApiResponseCommunity{Success: true, Message: "No pending communities found for user."}, nil
	}

	communitiesResponse := make([]*pb.Community, len(communities))
	for i, community := range communities {
		community, err := h.User_GetCommunityById(community.CommunityId, req.Value)
		if err != nil {
			zap.L().Error("Failed to get community by ID", zap.String("community_id", community.CommunityId), zap.Error(err))
			continue
		}
		communitiesResponse[i] = community
	}

	returnData, err := anypb.New(&pb.GetCommunitiesResponse{
		Communities: communitiesResponse,
	})

	if err != nil {
		return &pb.ApiResponseCommunity{
			Success: false,
			Message: "An error occurred: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponseCommunity{
		Success: true,
		Message: "Get all user pending approval communities successful.",
		Data:    returnData,
	}, nil
}
