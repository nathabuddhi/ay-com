package handlers

import (
	"context"
	"time"

	"github.com/nathabuddhi/ay-com/backend/service-community/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-community/proto/community"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandlers(db *gorm.DB) *Handler {
	return &Handler{
		DB: db,
	}
}

func (h *Handler) Community_GetAllCommunities(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	zap.L().Info("Getting all communities")

	var communities []models.Community
	err := h.DB.WithContext(ctx).Where("is_pending = ? AND is_rejected = ?", false, false).Find(&communities).Error
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
		getAllCommunitiesResponse.Communities[i], _ = h.User_GetCommunityById(community.CommunityId, req.Value)
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

	for _, category := range req.Categories {
		newCategory := models.CommunityCategoryRelation{
			CommunityId: req.CommunityId,
			Category:    category,
		}

		if err := h.DB.WithContext(ctx).Create(&newCategory).Error; err != nil {
			h.DB.WithContext(ctx).Delete(&community)
			return &pb.ApiResponseCommunity{Success: false, Message: "Failed to register category: " + err.Error()}, nil
		}
	}

	return &pb.ApiResponseCommunity{
		Success: true,
		Message: "Create community successful.",
		Data:    nil,
	}, nil
}

func (h *Handler) Community_GetAllUserCommunities(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	zap.L().Info("Getting all communities for user", zap.String("user_id", req.Value))

	var membership []models.CommunityMember
	err := h.DB.WithContext(ctx).Where("user_id = ?", req.Value).Find(&membership).Error
	if err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: err.Error()}, nil
	}

	if len(membership) == 0 {
		return &pb.ApiResponseCommunity{Success: true, Message: "No communities found for user."}, nil
	}

	communities := make([]*pb.Community, len(membership))
	for i, member := range membership {
		community, err := h.User_GetCommunityById(member.CommunityId, req.Value)
		if err != nil {
			zap.L().Error("Failed to get community by ID", zap.String("community_id", member.CommunityId), zap.Error(err))
			continue
		}
		communities[i] = community
	}

	returnData, err := anypb.New(&pb.GetCommunitiesResponse{
		Communities: communities,
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
		Message: "Get all user communities successful.",
		Data:    returnData,
	}, nil
}

func (h *Handler) Community_GetAllUserPendingCommunities(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	zap.L().Info("Getting all pending communities for user", zap.String("user_id", req.Value))

	var membership []models.CommunityJoinRequest
	err := h.DB.WithContext(ctx).Where("user_id = ? AND status = ?", req.Value, "pending").Find(&membership).Error
	if err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: err.Error()}, nil
	}

	if len(membership) == 0 {
		return &pb.ApiResponseCommunity{Success: true, Message: "No pending communities found for user."}, nil
	}

	communities := make([]*pb.Community, len(membership))
	for i, member := range membership {
		community, err := h.User_GetCommunityById(member.CommunityId, req.Value)
		if err != nil {
			zap.L().Error("Failed to get community by ID", zap.String("community_id", member.CommunityId), zap.Error(err))
			continue
		}
		communities[i] = community
	}

	returnData, err := anypb.New(&pb.GetCommunitiesResponse{
		Communities: communities,
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
		Message: "Get all user communities successful.",
		Data:    returnData,
	}, nil
}

func (h *Handler) Community_GetCommunityById(ctx context.Context, req *pb.GeneralCommunityRequest) (*pb.ApiResponseCommunity, error) {
	zap.L().Info("Getting community by ID", zap.String("community_id", req.CommunityId))

	community, err := h.User_GetCommunityById(req.CommunityId, req.UserId)
	if err != nil {
		return &pb.ApiResponseCommunity{Success: false, Message: "Failed to get community: " + err.Error()}, nil
	}

	returnData, err := anypb.New(community)

	if err != nil {
		return &pb.ApiResponseCommunity{
			Success: false,
			Message: "An error occurred while creating response: " + err.Error(),
		}, nil
	}

	return &pb.ApiResponseCommunity{
		Success: true,
		Message: "Community retrieved successfully.",
		Data:    returnData,
	}, nil
}
