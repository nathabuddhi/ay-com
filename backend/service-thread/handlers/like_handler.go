package handlers

import (
	"context"
	"time"

	"github.com/nathabuddhi/ay-com/backend/service-thread/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-thread/proto/thread"
	"github.com/nathabuddhi/ay-com/backend/service-thread/rabbitmq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (h *Handler) IsUserLikingThread(threadId string, userId string) bool {
	err := h.DB.Where("thread_id = ? AND user_id = ?", threadId, userId).First(&models.ThreadLike{}).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		zap.L().Error("Error checking if user is liking thread", zap.Error(err))
		return false
	}
	return true
}

func (h *Handler) Thread_ToggleLike(ctx context.Context, req *pb.GeneralThreadRequest) (*pb.ApiResponseThread, error) {
	zap.L().Info("Toggling like from ", zap.String("user_id", req.UserId))

	if req.ThreadId == "" || req.UserId == "" {
		return &pb.ApiResponseThread{Success: false, Message: "Invalid request parameters."}, nil
	}

	var threadLike models.ThreadLike
	err := h.DB.WithContext(ctx).Where("thread_id = ? AND user_id = ?", req.ThreadId, req.UserId).First(&threadLike).Error
	if err == nil {
		if err := h.DB.WithContext(ctx).Delete(&threadLike).Error; err != nil {
			return &pb.ApiResponseThread{Success: false, Message: "Failed to unlike thread."}, nil
		}
		return &pb.ApiResponseThread{Success: true, Message: "Thread unliked successfully."}, nil
	} else if err == gorm.ErrRecordNotFound {
		newLike := models.ThreadLike{
			ThreadId:  req.ThreadId,
			UserId:    req.UserId,
			CreatedAt: time.Now(),
		}
		if err := h.DB.WithContext(ctx).Create(&newLike).Error; err != nil {
			return &pb.ApiResponseThread{Success: false, Message: "Failed to like thread."}, nil
		}

		var thread models.Thread
		h.DB.WithContext(ctx).Where("thread_id = ?", req.ThreadId).First(&thread)
		rabbitmq.PublishSendNotification("like", thread.UserId, "", "Someone liked your thread.", req.UserId+" liked <a href='/thread/"+req.ThreadId+"'>your thread</a>", req.UserId)

		return &pb.ApiResponseThread{Success: true, Message: "Thread liked successfully."}, nil
	}

	return &pb.ApiResponseThread{Success: false, Message: "Unknown error occurred."}, nil
}

func (h *Handler) getLikeCount(ctx context.Context, threadId string) int {
	var count int64
	h.DB.WithContext(ctx).Model(&models.ThreadLike{}).Where("thread_id = ?", threadId).Count(&count)
	return int(count)
}
