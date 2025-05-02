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
	DeletedAt       *time.Time `gorm:"index"`
}

type Media struct {
	Id        string `gorm:"type:char(36);primaryKey"`
	ThreadId  string `gorm:"type:char(36);not null;index"`
	MediaURL  string `gorm:"type:text;not null"`
	MediaType string `gorm:"type:varchar(20);not null"`
}

type PollOption struct {
	Id        string `gorm:"type:char(36);primaryKey"`
	ThreadId  string `gorm:"type:char(36);not null;index"`
	Option    string `gorm:"type:varchar(255);not null"`
	VoteCount int    `gorm:"default:0"`
}

type ThreadLike struct {
	Id        string `gorm:"type:char(36);primaryKey"`
	UserId    string `gorm:"type:char(36);not null;index"`
	ThreadId  string `gorm:"type:char(36);not null;index"`
	CreatedAt time.Time
}

type ThreadRepost struct {
	Id        string `gorm:"type:char(36);primaryKey"`
	UserId    string `gorm:"type:char(36);not null;index"`
	ThreadId  string `gorm:"type:char(36);not null;index"`
	Text      string `gorm:"type:text"`
	CreatedAt time.Time
}

type ThreadReply struct {
	Id        string  `gorm:"type:char(36);primaryKey"`
	ThreadId  string  `gorm:"type:char(36);not null;index"`
	ReplyToId *string `gorm:"type:char(36);"`
	UserID    string  `gorm:"type:char(36);not null"`
	Content   string  `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}
