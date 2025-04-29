package models

type VerificationCode struct {
	Email string `gorm:"primaryKey"`
	Code  string
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
