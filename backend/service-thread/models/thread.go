package models

import (
	"time"
)

type Thread struct {
	ThreadId        string  `gorm:"type:char(36);primaryKey"`
	UserId          string  `gorm:"type:char(36);not null"`
	CommunityId     *string `gorm:"type:char(36);"`
	Content         string  `gorm:"type:text;not null"`
	Category        string  `gorm:"type:varchar(100);not null"`
	IsPoll          bool    `gorm:"default:false"`
	IsPrivate       bool    `gorm:"default:false"`
	IsScheduled     bool    `gorm:"default:false"`
	ScheduledAt     *time.Time
	IsAdvertisement bool   `gorm:"default:false"`
	HasMedia        bool   `gorm:"default:false"`
	ReplyPermission string `gorm:"type:varchar(20);default:'everyone'"`
	Pinned          bool   `gorm:"default:false"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Media struct {
	ThreadId  string `gorm:"type:char(36);not null;index;primaryKey"`
	MediaURL  string `gorm:"type:text;not null;primaryKey"`
	MediaType string `gorm:"type:varchar(20);not null"`
}

type PollOption struct {
	ThreadId string `gorm:"type:char(36);not null;primaryKey"`
	Option   string `gorm:"type:varchar(255);not null;primaryKey"`
}

type PollVote struct {
	ThreadId string `gorm:"type:char(36);not null;primaryKey"`
	Option   string `gorm:"type:varchar(255);not null"`
	UserId   string `gorm:"type:char(36);not null;primaryKey"`
}

type ThreadLike struct {
	UserId    string `gorm:"type:char(36);not null;index;primaryKey"`
	ThreadId  string `gorm:"type:char(36);not null;index;primaryKey"`
	CreatedAt time.Time
}

type ThreadBookmark struct {
	UserId    string `gorm:"type:char(36);not null;index;primaryKey"`
	ThreadId  string `gorm:"type:char(36);not null;index;primaryKey"`
	CreatedAt time.Time
}

type ThreadRepost struct {
	UserId    string `gorm:"type:char(36);not null;index;primaryKey"`
	ThreadId  string `gorm:"type:char(36);not null;index;primaryKey"`
	Text      string `gorm:"type:text"`
	CreatedAt time.Time
}

type ThreadReply struct {
	Id        string `gorm:"type:char(36);primaryKey"`
	ThreadId  string `gorm:"type:char(36);not null;index"`
	UserID    string `gorm:"type:char(36);not null"`
	Content   string `gorm:"type:text;not null"`
	IsPinned  bool   `gorm:"default:false"`
	CreatedAt time.Time
}
