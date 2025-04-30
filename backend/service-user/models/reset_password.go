package models

import "time"

type ResetPasswordCodes struct {
	email  string `gorm:"primaryKey"`
	code   string
	expiry time.Time
}
