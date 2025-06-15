package handlers

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-community/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-community/proto/community"
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
