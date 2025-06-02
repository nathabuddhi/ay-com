package handlers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nathabuddhi/ay-com/backend/service-thread/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-thread/proto/thread"
	"github.com/nathabuddhi/ay-com/backend/service-thread/rabbitmq"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
)

func (h *Handler) getReplyCount(ctx context.Context, threadId string) int {
	var count int64
	h.DB.WithContext(ctx).Model(&models.ThreadReply{}).Where("thread_id = ?", threadId).Count(&count)
	return int(count)
}

func (h *Handler) Thread_ReplyThread(ctx context.Context, req *pb.ReplyThreadRequest) (*pb.ApiResponseThread, error) {
	zap.L().Info("User " + req.UserId + " is replying to thread " + req.ThreadId)
	if req.Content == "" || req.ThreadId == "" || req.UserId == "" {
		return &pb.ApiResponseThread{
			Success: false,
			Message: "Incomplete request parameters.",
		}, nil
	}

	generatedId := uuid.New().String()

	reply := &models.ThreadReply{
		Id:        generatedId,
		ThreadId:  req.ThreadId,
		UserID:    req.UserId,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	if err := h.DB.WithContext(ctx).Create(reply).Error; err != nil {
		return &pb.ApiResponseThread{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	returnReply := &pb.ThreadReply{
		Id:        reply.Id,
		UserId:    reply.UserID,
		Content:   reply.Content,
		IsPinned:  false,
		Timestamp: reply.Content,
	}
	returnData, err := anypb.New(returnReply)

	if err != nil {
		return &pb.ApiResponseThread{
			Success: false,
			Message: err.Error(),
		}, err
	}

	rabbitmq.PublishDeleteRedis("getthread/" + req.ThreadId)

	return &pb.ApiResponseThread{
		Success: true,
		Message: "Reply added successfully.",
		Data:    returnData,
	}, nil
}
