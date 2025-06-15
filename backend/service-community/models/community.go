package models

import "time"

type Community struct {
	CommunityId   string `gorm:"primaryKey"`
	CommunityName string
	CreatorId     string
	Description   string `gorm:"type:text"`
	Rules         string `gorm:"type:text"`
	IconImage     string
	BannerImage   string
	IsPending     bool
	CreatedAt     time.Time
}

type CommunityCategory struct {
	Category string `gorm:"unique"`
}

type CommunityCategoryRelation struct {
	CommunityId string `gorm:"primaryKey"`
	Category    string `gorm:"primaryKey"`
}

type CommunityMember struct {
	CommunityId string `gorm:"primaryKey"`
	UserId      string `gorm:"primaryKey"`
	Role        string
	JoinedAt    time.Time
}

type CommunityJoinRequest struct {
	RequestId   string `gorm:"primaryKey"`
	CommunityId string
	UserId      string
	Status      string
	RequestedAt time.Time
}
