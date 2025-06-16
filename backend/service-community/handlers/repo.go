package handlers

import (
	"github.com/nathabuddhi/ay-com/backend/service-community/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-community/proto/community"
	"go.uber.org/zap"
)

func (h *Handler) GetMemberCount(communityId string) int32 {

	var memberCount int64
	err := h.DB.Where("community_id = ?", communityId).Model(&models.CommunityMember{}).Count(&memberCount).Error
	if err != nil {
		return 1
	}
	zap.L().Info("GetMemberCount",
		zap.String("community_id", communityId),
		zap.Int64("member_count", memberCount),
	)
	return int32(memberCount)
}

func (h *Handler) Admin_GetCommunityById(communityId string) (*pb.Community, error) {
	var community models.Community
	err := h.DB.Where("community_id = ? AND is_rejected = ?", communityId, false).First(&community).Error
	if err != nil {
		return nil, err
	}

	return &pb.Community{
		CommunityId:   community.CommunityId,
		CommunityName: community.CommunityName,
		Description:   community.Description,
		IconImage:     community.IconImage,
		BannerImage:   community.BannerImage,
		CreatorId:     community.CreatorId,
		CreatedAt:     community.CreatedAt.String(),
		Categories:    h.GetCommunityCategories(communityId),
		MemberCount:   h.GetMemberCount(community.CommunityId),
	}, nil
}

func (h *Handler) User_GetCommunityById(communityId string, requesterId string) (*pb.Community, error) {
	var community models.Community
	err := h.DB.Where("community_id = ? AND is_rejected = ?", communityId, false).First(&community).Error
	if err != nil {
		return nil, err
	}

	return &pb.Community{
		CommunityId:   community.CommunityId,
		CommunityName: community.CommunityName,
		Description:   community.Description,
		IconImage:     community.IconImage,
		BannerImage:   community.BannerImage,
		CreatorId:     community.CreatorId,
		CreatedAt:     community.CreatedAt.String(),
		Rules:         community.Rules,
		Categories:    h.GetCommunityCategories(communityId),
		MemberCount:   h.GetMemberCount(community.CommunityId),
		Role:          h.GetUserRole(requesterId, communityId),
		IsPending:     h.IsUserPendingJoin(requesterId, communityId),
	}, nil
}
