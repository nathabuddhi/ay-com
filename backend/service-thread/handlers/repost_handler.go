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

func (h *Handler) getRepostCount(ctx context.Context, threadId string) int {
	var count int64
	h.DB.WithContext(ctx).Model(&models.ThreadRepost{}).Where("thread_id = ?", threadId).Count(&count)
	return int(count)
}

func (h *Handler) IsUserRepostingThread(threadId string, userId string) bool {
	err := h.DB.Where("thread_id = ? AND user_id = ?", threadId, userId).First(&models.ThreadRepost{}).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		zap.L().Error("Error checking if user is reposting thread", zap.Error(err))
		return false
	}
	return true
}

func (h *Handler) Thread_ToggleRepost(ctx context.Context, req *pb.RepostRequest) (*pb.ApiResponseThread, error) {
	zap.L().Info("Toggling report from ", zap.String("user_id", req.UserId))

	if req.ThreadId == "" || req.UserId == "" {
		return &pb.ApiResponseThread{Success: false, Message: "Invalid request parameters."}, nil
	}

	var threadRepost models.ThreadRepost
	err := h.DB.WithContext(ctx).Where("thread_id = ? AND user_id = ?", req.ThreadId, req.UserId).First(&threadRepost).Error
	if err == nil {
		if err := h.DB.WithContext(ctx).Delete(&threadRepost).Error; err != nil {
			return &pb.ApiResponseThread{Success: false, Message: "Failed to unrepost thread."}, nil
		}
		return &pb.ApiResponseThread{Success: true, Message: "Thread unreposted successfully."}, nil
	} else if err == gorm.ErrRecordNotFound {
		var thread models.Thread
		h.DB.WithContext(ctx).Where("thread_id = ?", req.ThreadId).First(&thread)
		if thread.ThreadId == "" {
			return &pb.ApiResponseThread{Success: false, Message: "Thread not found."}, nil
		} else if thread.IsPrivate {
			return &pb.ApiResponseThread{Success: false, Message: "Cannot repost a private account's thread."}, nil
		} else if thread.IsAdvertisement {
			return &pb.ApiResponseThread{Success: false, Message: "Cannot repost an advertisement thread."}, nil
		} else if thread.UserId == req.UserId {
			return &pb.ApiResponseThread{Success: false, Message: "Cannot repost your own thread."}, nil
		}

		newRepost := models.ThreadRepost{
			ThreadId:  req.ThreadId,
			UserId:    req.UserId,
			Text:      req.Text,
			CreatedAt: time.Now(),
		}
		if err := h.DB.WithContext(ctx).Create(&newRepost).Error; err != nil {
			return &pb.ApiResponseThread{Success: false, Message: "Failed to repost thread."}, nil
		}
		return &pb.ApiResponseThread{Success: true, Message: "Thread reposted successfully."}, nil
	}

	return &pb.ApiResponseThread{Success: false, Message: "Unknown error occurred."}, nil
}

func (h *Handler) Thread_GetRepostedThreads(ctx context.Context, req *pb.StringThread) (*pb.ApiResponseThread, error) {
	zap.L().Info("Getting reposted threads from ", zap.String("user_id", req.Value))

	var reposts []models.ThreadRepost
	err := h.DB.WithContext(ctx).
		Where("user_id = ?", req.Value).
		Order("created_at desc").
		Find(&reposts).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.ApiResponseThread{Success: true, Message: "No reposts found."}, nil
		}
		return &pb.ApiResponseThread{Success: false, Message: err.Error()}, nil
	}

	threadIDs := make([]string, len(reposts))
	for i, bm := range reposts {
		threadIDs[i] = bm.ThreadId
	}

	var threads []models.Thread
	err = h.DB.WithContext(ctx).
		Where("thread_id IN ?", threadIDs).
		Find(&threads).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.ApiResponseThread{Success: true, Message: "No reposts found."}, nil
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
		Message: "Get reposts successful.",
		Data:    returnData,
	}, nil
}
