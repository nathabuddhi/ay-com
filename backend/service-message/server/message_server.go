package server

import (
	"github.com/nathabuddhi/ay-com/backend/service-message/handlers"
	pb "github.com/nathabuddhi/ay-com/backend/service-message/proto/message"
	"gorm.io/gorm"
)

type MessageServer struct {
	pb.UnimplementedMessageServiceServer
	Handler *handlers.Handler
}

func NewCommunityServer(db *gorm.DB) *MessageServer {
	return &MessageServer{
		Handler: handlers.NewHandler(db),
	}
}
