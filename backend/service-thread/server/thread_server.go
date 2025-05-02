package server

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-thread/handlers"
	pb "github.com/nathabuddhi/ay-com/backend/service-thread/proto/thread"
	"gorm.io/gorm"
)

type UserServer struct {
	pb.UnimplementedThreadServiceServer
	Handler *handlers.Handler
}

func NewThreadServer(db *gorm.DB) *UserServer {
	return &UserServer{
		Handler: handlers.NewHandler(db),
	}
}

func (s *UserServer) Thread_GetAllThreads(ctx context.Context, req *pb.GetAllThreadsRequest) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_GetAllThreads(ctx, req)
}
