package models

import "time"

type UserSetting struct {
	UserId    string `gorm:"primaryKey"`
	FontSize  string
	FontColor string
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
