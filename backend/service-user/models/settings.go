package models

import "time"

type UserSetting struct {
	UserId          string `gorm:"primaryKey"`
	FontSize        string
	FontColor       string
	NotifLike       bool
	NotifRepost     bool
	NotifFollow     bool
	NotifMention    bool
	NotifCommunity  bool
	NotifNewsletter bool
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}
