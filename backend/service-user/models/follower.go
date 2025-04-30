package models

import (
	"time"
)

type UserFollower struct {
	id          string `gorm:"primaryKey"`
	uesr_id     string
	follower_id string
	created_at  time.Time `gorm:"autoCreateTime"`
}
