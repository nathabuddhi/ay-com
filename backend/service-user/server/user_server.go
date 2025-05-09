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

func (s *UserServer) User_Login(ctx context.Context, req *pb.LoginRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_Login(ctx, req)
}

func (s *UserServer) User_Register(ctx context.Context, req *pb.RegisterRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_Register(ctx, req)
}

func (s *UserServer) User_GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_GetProfile(ctx, req)
}

func (s *UserServer) User_RequestVerificationCode(ctx context.Context, req *pb.VerificationRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_RequestVerificationCode(ctx, req)
}

func (s *UserServer) User_ValidateVerificationCode(ctx context.Context, req *pb.ValidateCodeRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_ValidateVerificationCode(ctx, req)
}

func (s *UserServer) User_ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_ChangePassword(ctx, req)
}

func (s *UserServer) User_GetSecurityQuestion(ctx context.Context, req *pb.GetSecurityQuestionRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_GetSecurityQuestion(ctx, req)
}

func (s *UserServer) User_ValidateSecurityAnswer(ctx context.Context, req *pb.ValidateSecurityAnswerRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_ValidateSecurityAnswer(ctx, req)
}

func (s *UserServer) User_ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_ResetPassword(ctx, req)
}

func (s *UserServer) User_FollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_FollowUser(ctx, req)
}

func (s *UserServer) User_BlockUser(ctx context.Context, req *pb.BlockUserRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_BlockUser(ctx, req)
}

func (s *UserServer) User_UnFollowUser(ctx context.Context, req *pb.UnFollowUserRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_UnFollowUser(ctx, req)
}

func (s *UserServer) User_UnBlockUser(ctx context.Context, req *pb.UnBlockUserRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_UnBlockUser(ctx, req)
}

func (s *UserServer) User_GetSettings(ctx context.Context, req *pb.GetSettingsRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_GetSettings(ctx, req)
}

func (s *UserServer) User_UpdateSettings(ctx context.Context, req *pb.UpdateSettingsRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_UpdateSettings(ctx, req)
}

func (s *UserServer) User_GetAllFollowers(ctx context.Context, req *pb.GetAllFollowersRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_GetAllFollowers(ctx, req)
}

func (s *UserServer) User_GetAllFollowing(ctx context.Context, req *pb.GetAllFollowingRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_GetAllFollowing(ctx, req)
}

func (s *UserServer) User_GetAllVerifyAccountRequest(ctx context.Context, req *pb.GetAllVerifyAccountRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_GetAllVerifyAccountRequest(ctx, req)
}

func (s *UserServer) User_SubmitVerifyAccountRequest(ctx context.Context, req *pb.SubmitVerifyAccountRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_SubmitVerifyAccountRequest(ctx, req)
}

func (s *UserServer) User_DeactivateAccount(ctx context.Context, req *pb.DeactivateAccountRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_DeactivateAccount(ctx, req)
}

func (s *UserServer) User_UpdateProfile(ctx context.Context, req *pb.UpdateUserProfileRequest) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_UpdateProfile(ctx, req)
}

func (s *UserServer) User_CheckToken(ctx context.Context, req *pb.StringUser) (*pb.BoolUser, error) {
	return s.Handlers.User_CheckToken(ctx, req)
}

func (s *UserServer) User_GetSelfProfile(ctx context.Context, req *pb.StringUser) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_GetSelfProfile(ctx, req)
}

func (s *UserServer) User_GetAllBlocked(ctx context.Context, req *pb.StringUser) (*pb.ApiResponseUser, error) {
	return s.Handlers.User_GetAllBlocked(ctx, req)
}

// func (s *UserServer) User_(ctx context.Context, req *pb.) (*pb.ApiResponseUser, error) {
// 	return s.Handlers.User_(ctx, req)
// }
