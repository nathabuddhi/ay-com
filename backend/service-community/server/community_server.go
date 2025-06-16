package server

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-community/handlers"
	pb "github.com/nathabuddhi/ay-com/backend/service-community/proto/community"
)

type CommunityServer struct {
	pb.UnimplementedCommunityServiceServer
	Handler *handlers.Handler
}

func NewCommunityServer(newHandler *handlers.Handler) *CommunityServer {
	return &CommunityServer{
		Handler: newHandler,
	}
}

func (s *CommunityServer) Community_GetAllCommunities(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	return s.Handler.Community_GetAllCommunities(ctx, req)
}

func (s *CommunityServer) Community_GetUserCommunities(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	return s.Handler.Community_GetAllUserCommunities(ctx, req)
}

func (s *CommunityServer) Community_GetUserPendingCommunities(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	return s.Handler.Community_GetAllUserPendingCommunities(ctx, req)
}

func (s *CommunityServer) Community_GetUserPendingApprovalCommunities(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	return s.Handler.Community_GetAllUserPendingApprovalCommunities(ctx, req)
}

func (s *CommunityServer) Community_CreateCommunity(ctx context.Context, req *pb.CreateCommunityRequest) (*pb.ApiResponseCommunity, error) {
	return s.Handler.Community_CreateCommunity(ctx, req)
}

func (s *CommunityServer) Community_GetCommunityById(ctx context.Context, req *pb.GeneralCommunityRequest) (*pb.ApiResponseCommunity, error) {
	return s.Handler.Community_GetCommunityById(ctx, req)
}

func (s *CommunityServer) Community_GetCategories(ctx context.Context, req *pb.StringCommunity) (*pb.GetCategoriesResponse, error) {
	return s.Handler.Community_GetCategories(ctx, req)
}

func (s *CommunityServer) Community_GetCommunityMembers(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	return s.Handler.Community_GetCommunityMembers(ctx, req)
}

func (s *CommunityServer) Community_JoinCommunity(ctx context.Context, req *pb.GeneralCommunityRequest) (*pb.ApiResponseCommunity, error) {
	return s.Handler.Community_JoinCommunity(ctx, req)
}

// func (s *CommunityServer) Community_LeaveCommunity(ctx context.Context, req *pb.GeneralCommunityRequest) (*pb.ApiResponseCommunity, error) {
// 	return s.Handler.Community_LeaveCommunity(ctx, req)
// }

func (s *CommunityServer) Community_ApproveJoinRequest(ctx context.Context, req *pb.GeneralModeratorRequest) (*pb.ApiResponseCommunity, error) {
	return s.Handler.Community_ApproveJoinRequest(ctx, req)
}

func (s *CommunityServer) Community_DenyJoinRequest(ctx context.Context, req *pb.GeneralModeratorRequest) (*pb.ApiResponseCommunity, error) {
	return s.Handler.Community_DenyJoinRequest(ctx, req)
}

func (s *CommunityServer) Community_GetJoinRequests(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	return s.Handler.Community_GetPendingUsers(ctx, req)
}

// ADMIN
func (s *CommunityServer) Admin_GetAllCommunityRequests(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	return s.Handler.Admin_GetAllCommunityRequests(ctx, req)
}

func (s *CommunityServer) Admin_ApproveCommunity(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	return s.Handler.Admin_ApproveCommunity(ctx, req)
}

func (s *CommunityServer) Admin_RejectCommunity(ctx context.Context, req *pb.RejectCommunityRequest) (*pb.ApiResponseCommunity, error) {
	return s.Handler.Admin_RejectCommunity(ctx, req)
}

func (s *CommunityServer) Admin_AddCategory(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	return s.Handler.Admin_AddCategory(ctx, req)
}

func (s *CommunityServer) Admin_DeleteCategory(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	return s.Handler.Admin_DeleteCategory(ctx, req)
}
