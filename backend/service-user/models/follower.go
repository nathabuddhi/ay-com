package models

import (
	"time"
)

type UserFollowing struct {
	UserId     string    `gorm:"primaryKey"`
	FollowedId string    `gorm:"primaryKey"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

type BlockedUsers struct {
	UserId    string    `gorm:"primaryKey"`
	BlockedId string    `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
