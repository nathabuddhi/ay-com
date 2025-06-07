package handlers

import (
	"context"
	"fmt"

	"github.com/nathabuddhi/ay-com/backend/service-thread/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-thread/proto/thread"
	"go.uber.org/zap"
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
