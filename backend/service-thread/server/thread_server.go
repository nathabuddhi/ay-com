package server

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-thread/handlers"
	pb "github.com/nathabuddhi/ay-com/backend/service-thread/proto/thread"
	"gorm.io/gorm"
)

type ThreadServer struct {
	pb.UnimplementedThreadServiceServer
	Handler *handlers.Handler
}

func NewThreadServer(db *gorm.DB) *ThreadServer {
	return &ThreadServer{
		Handler: handlers.NewHandler(db),
	}
}

func (s *ThreadServer) Thread_GetAllThreads(ctx context.Context, req *pb.GetAllThreadsRequest) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_GetAllThreads(ctx, req)
}

func (s *ThreadServer) Thread_GetFollowingThreads(ctx context.Context, req *pb.GetFollowingThreadRequest) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_GetFollowingThreads(ctx, req)
}

func (s *ThreadServer) Thread_GetThreadById(ctx context.Context, req *pb.GeneralThreadRequest) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_GetThreadById(ctx, req)
}

func (s *ThreadServer) Thread_CreateThread(ctx context.Context, req *pb.PostThread) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_CreateThread(ctx, req)
}

// func (s *ThreadServer) Thread_SearchThreads(ctx context.Context, req *pb.GetAllThreadsRequest) (*pb.ApiResponseThread, error) {
// 	return s.Handler.Thread_SearchThreads(ctx, req)
// }

func (s *ThreadServer) Thread_DeleteThread(ctx context.Context, req *pb.GeneralThreadRequest) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_DeleteThread(ctx, req)
}

func (s *ThreadServer) Thread_TogglePinThread(ctx context.Context, req *pb.GeneralThreadRequest) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_TogglePinThread(ctx, req)
}

func (s *ThreadServer) Thread_VoteThread(ctx context.Context, req *pb.SubmitVote) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_VoteThread(ctx, req)
}

func (s *ThreadServer) Thread_ToggleLike(ctx context.Context, req *pb.GeneralThreadRequest) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_ToggleLike(ctx, req)
}

func (s *ThreadServer) Thread_ToggleBookmark(ctx context.Context, req *pb.GeneralThreadRequest) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_ToggleBookmark(ctx, req)
}

func (s *ThreadServer) Thread_ToggleRepost(ctx context.Context, req *pb.RepostRequest) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_ToggleRepost(ctx, req)
}

func (s *ThreadServer) Thread_GetBookmarkedThreads(ctx context.Context, req *pb.StringThread) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_GetBookmarkedThreads(ctx, req)
}

func (s *ThreadServer) Thread_GetRepostedThreads(ctx context.Context, req *pb.StringThread) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_GetRepostedThreads(ctx, req)
}

func (s *ThreadServer) Thread_GetReplyPermission(ctx context.Context, req *pb.StringThread) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_GetReplyPermission(ctx, req)
}

func (s *ThreadServer) Thread_GetUserThreads(ctx context.Context, req *pb.UserToUserRequest) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_GetUserThreads(ctx, req)
}

func (s *ThreadServer) Thread_GetUserLikedThreads(ctx context.Context, req *pb.StringThread) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_GetUserLikedThreads(ctx, req)
}

func (s *ThreadServer) Thread_GetUserMediaThreads(ctx context.Context, req *pb.UserToUserRequest) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_GetUserMediaThreads(ctx, req)
}

func (s *ThreadServer) Thread_GetUserReplies(ctx context.Context, req *pb.UserToUserRequest) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_GetUserReplies(ctx, req)
}

func (s *ThreadServer) Thread_GetTrendingHashtags(ctx context.Context, req *pb.StringThread) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_GetTrendingHashtags(ctx, req)
}

func (s *ThreadServer) Thread_GetThreadCategories(ctx context.Context, req *pb.StringThread) (*pb.ThreadCategories, error) {
	return s.Handler.Thread_GetThreadCategories(ctx, req)
}

func (s *ThreadServer) Thread_GetAdvertisementThreads(ctx context.Context, req *pb.GeneralThreadRequest) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_GetAdvertisementThreads(ctx, req)
}

func (s *ThreadServer) Thread_GetCommunityMediaThreads(ctx context.Context, req *pb.GeneralThreadRequest) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_GetCommunityMediaThreads(ctx, req)
}

func (s *ThreadServer) Thread_GetCommunityThreads(ctx context.Context, req *pb.GeneralThreadRequest) (*pb.ApiResponseThread, error) {
	return s.Handler.Thread_GetCommunityThreads(ctx, req)
}

// ADMIN
func (s *ThreadServer) Admin_DeleteThread(ctx context.Context, req *pb.StringThread) (*pb.ApiResponseThread, error) {
	return s.Handler.Admin_DeleteThread(ctx, req)
}

func (s *ThreadServer) Admin_AddCategory(ctx context.Context, req *pb.StringThread) (*pb.ApiResponseThread, error) {
	return s.Handler.Admin_AddCategory(ctx, req)
}

func (s *ThreadServer) Admin_DeleteCategory(ctx context.Context, req *pb.StringThread) (*pb.ApiResponseThread, error) {
	return s.Handler.Admin_DeleteCategory(ctx, req)
}
