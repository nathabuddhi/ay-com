package server

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-user/handlers"
	emailpb "github.com/nathabuddhi/ay-com/backend/service-user/proto/email"
	redispb "github.com/nathabuddhi/ay-com/backend/service-user/proto/redis"
	pb "github.com/nathabuddhi/ay-com/backend/service-user/proto/user"
	"gorm.io/gorm"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	Handlers *handlers.Handlers
}

func NewUserServer(db *gorm.DB, redisClient redispb.RedisServiceClient, emailClient emailpb.EmailServiceClient) *UserServer {
	return &UserServer{
		Handlers: handlers.NewHandlers(db, redisClient, emailClient),
	}
}

func (s *UserServer) User_Login(ctx context.Context, req *pb.LoginRequest) (*pb.ApiResponse, error) {
	return s.Handlers.Login(ctx, req)
}

func (s *UserServer) User_Register(ctx context.Context, req *pb.RegisterRequest) (*pb.ApiResponse, error) {
	return s.Handlers.Register(ctx, req)
}

func (s *UserServer) User_GetProfile(ctx context.Context, req *pb.GetUserRequest) (*pb.ApiResponse, error) {
	return s.Handlers.GetProfile(ctx, req)
}

func (s *UserServer) RequestVerificationCode(ctx context.Context, req *pb.VerificationRequest) (*pb.ApiResponse, error) {
	return s.Handlers.RequestVerificationCode(ctx, req)
}

func (s *UserServer) ValidateVerificationCode(ctx context.Context, req *pb.ValidateCodeRequest) (*pb.ApiResponse, error) {
	return s.Handlers.ValidateVerificationCode(ctx, req)
}
