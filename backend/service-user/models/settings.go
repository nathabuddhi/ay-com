package models

type UserSetting struct {
	UserId    string `gorm:"primaryKey"`
	FontSize  string
	FontColor string
}
