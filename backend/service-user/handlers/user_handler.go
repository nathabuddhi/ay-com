package handlers

import (
	emailpb "github.com/nathabuddhi/ay-com/backend/service-user/proto/email"
	redispb "github.com/nathabuddhi/ay-com/backend/service-user/proto/redis"
	"gorm.io/gorm"
)

type Handlers struct {
	DB          *gorm.DB
	RedisClient redispb.RedisServiceClient
	EmailClient emailpb.EmailServiceClient
}

func NewHandlers(db *gorm.DB, redisClient redispb.RedisServiceClient, emailClient emailpb.EmailServiceClient) *Handlers {
	return &Handlers{
		DB:          db,
		RedisClient: redisClient,
		EmailClient: emailClient,
	}
}
