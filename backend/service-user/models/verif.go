package models

import "time"

type VerificationCode struct {
	Email  string `gorm:"primaryKey"`
	Code   string
	Expiry time.Time
}

type ValidateVerificationCodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type PublishNotificationEmailRequest struct {
	Email   string
	Subject string
	Body    string
}
