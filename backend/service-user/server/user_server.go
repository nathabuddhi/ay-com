package server

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-user/handlers"
	pb "github.com/nathabuddhi/ay-com/backend/service-user/proto/user"
	"gorm.io/gorm"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	Handlers *handlers.Handlers
}

func NewUserServer(db *gorm.DB) *UserServer {
	return &UserServer{
		Handlers: handlers.NewHandlers(db),
	}
}

func (s *UserServer) User_Login(ctx context.Context, req *pb.LoginRequest) (*pb.ApiResponse, error) {
	return s.Handlers.User_Login(ctx, req)
}

func (s *UserServer) User_Register(ctx context.Context, req *pb.RegisterRequest) (*pb.ApiResponse, error) {
	return s.Handlers.User_Register(ctx, req)
}

func (s *UserServer) User_GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.ApiResponse, error) {
	return s.Handlers.User_GetProfile(ctx, req)
}

func (s *UserServer) User_RequestVerificationCode(ctx context.Context, req *pb.VerificationRequest) (*pb.ApiResponse, error) {
	return s.Handlers.User_RequestVerificationCode(ctx, req)
}

func (s *UserServer) User_ValidateVerificationCode(ctx context.Context, req *pb.ValidateCodeRequest) (*pb.ApiResponse, error) {
	return s.Handlers.User_ValidateVerificationCode(ctx, req)
}

func (s *UserServer) User_ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ApiResponse, error) {
	return s.Handlers.User_ChangePassword(ctx, req)
}

func (s *UserServer) User_GetSecurityQuestion(ctx context.Context, req *pb.GetSecurityQuestionRequest) (*pb.ApiResponse, error) {
	return s.Handlers.User_GetSecurityQuestion(ctx, req)
}

func (s *UserServer) User_ValidateSecurityAnswer(ctx context.Context, req *pb.ValidateSecurityAnswerRequest) (*pb.ApiResponse, error) {
	return s.Handlers.User_ValidateSecurityAnswer(ctx, req)
}

func (s *UserServer) User_ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ApiResponse, error) {
	return s.Handlers.User_ResetPassword(ctx, req)
}

func (s *UserServer) User_FollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.ApiResponse, error) {
	return s.Handlers.User_FollowUser(ctx, req)
}

func (s *UserServer) User_BlockUser(ctx context.Context, req *pb.BlockUserRequest) (*pb.ApiResponse, error) {
	return s.Handlers.User_BlockUser(ctx, req)
}

func (s *UserServer) User_UnFollowUser(ctx context.Context, req *pb.UnFollowUserRequest) (*pb.ApiResponse, error) {
	return s.Handlers.User_UnFollowUser(ctx, req)
}

func (s *UserServer) User_UnBlockUser(ctx context.Context, req *pb.UnBlockUserRequest) (*pb.ApiResponse, error) {
	return s.Handlers.User_UnBlockUser(ctx, req)
}

func (s *UserServer) User_GetSettings(ctx context.Context, req *pb.GetSettingsRequest) (*pb.ApiResponse, error) {
	return s.Handlers.User_GetSettings(ctx, req)
}

func (s *UserServer) User_UpdateSettings(ctx context.Context, req *pb.UpdateSettingsRequest) (*pb.ApiResponse, error) {
	return s.Handlers.User_UpdateSettings(ctx, req)
}

// func (s *UserServer) User_(ctx context.Context, req *pb.) (*pb.ApiResponse, error) {
// 	return s.Handlers.User_(ctx, req)
// }
