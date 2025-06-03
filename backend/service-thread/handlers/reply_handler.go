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
	"gorm.io/gorm"
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

func (h *Handler) Thread_TogglePinReply(ctx context.Context, req *pb.GeneralThreadRequest) (*pb.ApiResponseThread, error) {
	zap.L().Info("User "+req.UserId+" is toggling pin for reply", zap.String("reply_id", req.ThreadId))

	if req.ThreadId == "" || req.UserId == "" {
		return &pb.ApiResponseThread{Success: false, Message: "Invalid request parameters."}, nil
	}

	var threadReply models.ThreadReply
	err := h.DB.WithContext(ctx).Where("id = ?", req.ThreadId).First(&threadReply).Error
	if err == nil {
		err := h.DB.WithContext(ctx).Model(&threadReply).Update("is_pinned", !threadReply.IsPinned).Error
		if err != nil {
			return &pb.ApiResponseThread{Success: false, Message: "Failed to toggle pin status."}, nil
		}
		rabbitmq.PublishDeleteRedis("getthread/" + threadReply.ThreadId)
		return &pb.ApiResponseThread{Success: true, Message: "Pin status toggled successfully."}, nil
	} else if err == gorm.ErrRecordNotFound {
		return &pb.ApiResponseThread{Success: false, Message: "Reply not found."}, nil
	} else {
		return &pb.ApiResponseThread{Success: false, Message: err.Error()}, nil
	}
}

func (h *Handler) Thread_DeleteReply(ctx context.Context, req *pb.GeneralThreadRequest) (*pb.ApiResponseThread, error) {
	zap.L().Info("User "+req.UserId+" is deleting reply", zap.String("reply_id", req.ThreadId))

	if req.ThreadId == "" || req.UserId == "" {
		return &pb.ApiResponseThread{Success: false, Message: "Invalid request parameters."}, nil
	}

	var thread models.Thread
	err := h.DB.WithContext(ctx).Where("thread_id = ?", req.ThreadId).First(&thread).Error
	if err == nil {
		err := h.DB.WithContext(ctx).Delete(&thread).Error
		if err != nil {
			return &pb.ApiResponseThread{Success: false, Message: "Failed to delete thread."}, nil
		}
		rabbitmq.PublishDeleteRedis("getthread/" + thread.ThreadId)
		return &pb.ApiResponseThread{Success: true, Message: "Thread deleted successfully."}, nil
	} else {
		return &pb.ApiResponseThread{
			Success: false,
			Message: "Thread not found or you do not have permission to delete this thread.",
		}, nil
	}
}
