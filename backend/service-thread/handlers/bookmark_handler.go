package handlers

import (
	"context"
	"time"

	"github.com/nathabuddhi/ay-com/backend/service-thread/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-thread/proto/thread"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"
)

func (h *Handler) IsUserBookmarkingThread(threadId string, userId string) bool {
	err := h.DB.Where("thread_id = ? AND user_id = ?", threadId, userId).First(&models.ThreadBookmark{}).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		zap.L().Error("Error checking if user is bookmarking thread", zap.Error(err))
		return false
	}
	return true
}

func (h *Handler) Thread_ToggleBookmark(ctx context.Context, req *pb.GeneralThreadRequest) (*pb.ApiResponseThread, error) {
	zap.L().Info("Toggling bookmark from ", zap.String("user_id", req.UserId))

	if req.ThreadId == "" || req.UserId == "" {
		return &pb.ApiResponseThread{Success: false, Message: "Invalid request parameters."}, nil
	}

	var threadBookmark models.ThreadBookmark
	err := h.DB.WithContext(ctx).Where("thread_id = ? AND user_id = ?", req.ThreadId, req.UserId).First(&threadBookmark).Error
	if err == nil {
		if err := h.DB.WithContext(ctx).Delete(&threadBookmark).Error; err != nil {
			return &pb.ApiResponseThread{Success: false, Message: "Failed to unbookmark thread."}, nil
		}
		return &pb.ApiResponseThread{Success: true, Message: "Thread unbookmarked successfully."}, nil
	} else if err == gorm.ErrRecordNotFound {
		newBookmark := models.ThreadBookmark{
			ThreadId:  req.ThreadId,
			UserId:    req.UserId,
			CreatedAt: time.Now(),
		}
		if err := h.DB.WithContext(ctx).Create(&newBookmark).Error; err != nil {
			return &pb.ApiResponseThread{Success: false, Message: "Failed to bookmark thread."}, nil
		}
		return &pb.ApiResponseThread{Success: true, Message: "Thread bookmarked successfully."}, nil
	}

	return &pb.ApiResponseThread{Success: false, Message: "Unknown error occurred."}, nil
}

func (h *Handler) Thread_GetBookmarkedThreads(ctx context.Context, req *pb.StringThread) (*pb.ApiResponseThread, error) {
	zap.L().Info("Getting bookmarked threads from ", zap.String("user_id", req.Value))

	var bookmarks []models.ThreadBookmark
	err := h.DB.WithContext(ctx).
		Where("user_id = ?", req.Value).
		Order("created_at desc").
		Find(&bookmarks).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.ApiResponseThread{Success: true, Message: "No bookmarks found."}, nil
		}
		return &pb.ApiResponseThread{Success: false, Message: err.Error()}, nil
	}

	threadIDs := make([]string, len(bookmarks))
	for i, bm := range bookmarks {
		threadIDs[i] = bm.ThreadId
	}

	var threads []models.Thread
	err = h.DB.WithContext(ctx).
		Where("thread_id IN ?", threadIDs).
		Find(&threads).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.ApiResponseThread{Success: true, Message: "No bookmarks found."}, nil
		}
		return &pb.ApiResponseThread{Success: false, Message: err.Error()}, nil
	}

	getAllThreadResponse := &pb.GetThreadsResponse{
		Threads: make([]*pb.Thread, len(threads)),
	}

	for i, request := range threads {
		threadResponse, err := h.processThreadResponse(ctx, request, req.Value)
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
		Message: "Get bookmarks successful",
		Data:    returnData,
	}, nil
}
