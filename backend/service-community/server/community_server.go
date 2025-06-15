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

func (s *CommunityServer) Community_GetAllUserCommunities(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	return s.Handler.Admin_GetAllCommunityRequests(ctx, req)
}

func (s *CommunityServer) Community_CreateCommunity(ctx context.Context, req *pb.CreateCommunityRequest) (*pb.ApiResponseCommunity, error) {
	return s.Handler.Community_CreateCommunity(ctx, req)
}

// ADMIN
func (s *CommunityServer) Admin_GetAllCommunityRequests(ctx context.Context, req *pb.StringCommunity) (*pb.ApiResponseCommunity, error) {
	return s.Handler.Admin_GetAllCommunityRequests(ctx, req)
}
