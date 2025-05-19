package models

import (
	"time"
)

type User struct {
	UserId           string `gorm:"primaryKey"`
	Name             string
	Username         string `gorm:"uniqueIndex"`
	Email            string `gorm:"uniqueIndex"`
	Password         string
	Gender           string
	DateOfBirth      time.Time
	IsVerified       bool
	IsBanned         bool
	IsDeactivated    bool
	IsPrivate        bool
	Bio              *string
	SecurityQuestion string
	SecurityAnswer   string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	JoinedAt         time.Time
	GoogleAuthID     *string
}

type UserVerificationRequest struct {
	Id                 string `gorm:"primaryKey"`
	UserId             string
	IdentityCardNumber string `gorm:"type:text"`
	SelfieUrl          string `gorm:"type:text"`
	ReasonText         string `gorm:"type:text"`
	Status             string
	SubmittedAt        time.Time
}

type RefreshToken struct {
	UserId string `gorm:"primaryKey"`
	Token  string `gorm:"type:text"`
}
