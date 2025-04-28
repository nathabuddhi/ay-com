package models

import (
	"time"
)

type User struct {
	UserId           string `gorm:"type:uuid;primaryKey"`
	Name             string
	Username         string `gorm:"uniqueIndex"`
	Email            string `gorm:"uniqueIndex"`
	Password         string
	Gender           string
	DateOfBirth      time.Time
	ProfilePicture   *string
	Banner           *string
	IsVerified       bool
	IsBanned         bool
	IsDeactivated    bool
	IsPrivate        bool
	Bio              *string
	SecurityQuestion string
	SecurityAnswer   string
	WantsNewsletter  bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	JoinedAt         time.Time
	GoogleAuthID     *string
}
