package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/nathabuddhi/ay-com/backend/service-thread/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-thread/proto/thread"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"
)

func (h *Handler) getReplyCount(ctx context.Context, threadId string) int {
	var count int64
	h.DB.WithContext(ctx).Model(&models.ThreadReply{}).Where("reply_to_id = ?", threadId).Count(&count)
	return int(count)
}

func (h *Handler) Thread_ReplyThread(ctx context.Context, threadId string, replyToId string) error {
	zap.L().Info("Thread " + threadId + " is replying to thread " + replyToId)
	if threadId == "" || replyToId == "" {
		return fmt.Errorf("no thread or reply provided")
	}

	reply := &models.ThreadReply{
		ThreadId:  threadId,
		ReplyToId: replyToId,
	}

	if err := h.DB.WithContext(ctx).Create(reply).Error; err != nil {
		return err
	}

	return nil
}

func (h *Handler) Thread_GetReplyPermission(ctx context.Context, req *pb.StringThread) (*pb.ApiResponseThread, error) {
	if req.Value == "" {
		return &pb.ApiResponseThread{
			Success: false,
			Message: "No thread ID provided",
			Data:    nil,
		}, nil
	}

	var thread models.Thread
	if err := h.DB.WithContext(ctx).First(&thread, "id = ?", thread).Error; err != nil {
		return &pb.ApiResponseThread{
			Success: false,
			Message: "Thread not found!",
			Data:    nil,
		}, nil
	} else {
		return &pb.ApiResponseThread{
			Success: true,
			Message: thread.ReplyPermission,
			Data:    nil,
		}, nil
	}
}

func (h *Handler) Thread_GetUserReplies(ctx context.Context, req *pb.UserToUserRequest) (*pb.ApiResponseThread, error) {
	zap.L().Info("Getting user replies ", zap.String("user_id", req.UserId))

	var threads []models.Thread
	err := h.DB.WithContext(ctx).
		Where("is_scheduled = false OR scheduled_at < ?", time.Now()).
		Where("user_id = ? AND category = ?", req.UserId, "Reply").
		Find(&threads).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.ApiResponseThread{Success: true, Message: "No replies found."}, nil
		}
		return &pb.ApiResponseThread{Success: false, Message: err.Error()}, nil
	}

	getAllThreadResponse := &pb.GetThreadsResponse{
		Threads: make([]*pb.Thread, len(threads)),
	}

	for i, request := range threads {
		threadResponse, err := h.processThreadResponse(ctx, request, req.RequesterId)
		if err != nil {
			return &pb.ApiResponseThread{
				Success: false,
				Message: err.Error(),
			}, nil
		}
		getAllThreadResponse.Threads[i] = &threadResponse
	}

	returnData, err := anypb.New(getAllThreadResponse)
	if err != nil {
		return &pb.ApiResponseThread{
			Success: false,
			Message: "An error occured: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponseThread{
		Success: true,
		Message: "Get threads successful.",
		Data:    returnData,
	}, nil
}
