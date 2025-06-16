package models

type ChatRoom struct {
	ChatRoomId   string `gorm:"primaryKey;"`
	ChatRoomName string `gorm:"not null;"`
	ChatRoomType string `gorm:"not null;"`
}

type ChatRoomMember struct {
	ChatRoomId string `gorm:"primaryKey;"`
	UserId     string `gorm:"primaryKey;"`
}

type Message struct {
	MessageId  string `gorm:"primaryKey;"`
	ChatRoomId string `gorm:"not null;"`
	UserId     string `gorm:"not null;"`
	Content    string `gorm:"not null;"`
	Timestamp  int64  `gorm:"not null;"`
}

type OpenConversations struct {
	UserId     string `gorm:"primaryKey;"`
	ChatRoomId string `gorm:"primaryKey;"`
}
