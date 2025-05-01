package models

import "time"

type Notification struct {
	NotificationId string `gorm:"primaryKey"`
	UserId         string
	Title          string
	Content        string
	Read           bool
	From           string
	Timestamp      time.Time
}

type NotificationSetting struct {
	UserId          string `gorm:"primaryKey"`
	NotifLike       bool
	NotifRepost     bool
	NotifFollow     bool
	NotifMention    bool
	NotifCommunity  bool
	NotifNewsletter bool
}
