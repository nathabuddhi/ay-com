package handlers

import (
	"github.com/nathabuddhi/ay-com/backend/service-community/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-community/proto/community"
)

func (h *Handler) GetMemberCount(communityId string) int32 {

	var memberCount int64
	err := h.DB.Where("community_id = ?", communityId).Model(&models.CommunityMember{}).Count(&memberCount).Error
	if err != nil {
		return -1
	}

	return int32(memberCount)
}

func (h *Handler) GetCommunityById(communityId string) (*pb.Community, error) {
	var community models.Community
	err := h.DB.Where("community_id = ?", communityId).First(&community).Error
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
		MemberCount:   h.GetMemberCount(community.CommunityId),
	}, nil
}
