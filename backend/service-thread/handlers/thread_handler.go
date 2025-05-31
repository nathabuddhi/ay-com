package handlers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nathabuddhi/ay-com/backend/service-thread/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-thread/proto/thread"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		DB: db,
	}
}

func (h *Handler) getReplyCount(ctx context.Context, threadId string) int {
	var count int64
	h.DB.WithContext(ctx).Model(&models.ThreadReply{}).Where("thread_id = ?", threadId).Count(&count)
	return int(count)
}

func (h *Handler) getLikeCount(ctx context.Context, threadId string) int {
	var count int64
	h.DB.WithContext(ctx).Model(&models.ThreadLike{}).Where("thread_id = ?", threadId).Count(&count)
	return int(count)
}

func (h *Handler) getRepostCount(ctx context.Context, threadId string) int {
	var count int64
	h.DB.WithContext(ctx).Model(&models.ThreadRepost{}).Where("thread_id = ?", threadId).Count(&count)
	return int(count)
}

func (h *Handler) getMedia(ctx context.Context, threadId string) ([]*pb.Media, error) {
	var mediaList []models.Media
	err := h.DB.WithContext(ctx).Where("thread_id = ?", threadId).Find(&mediaList).Error
	if err != nil {
		return nil, err
	}

	var mediaResponse []*pb.Media
	for _, media := range mediaList {
		mediaResponse = append(mediaResponse, &pb.Media{
			MediaUrl:  media.MediaURL,
			MediaType: media.MediaType,
		})
	}

	return mediaResponse, nil
}

func (h *Handler) getPollOptions(ctx context.Context, threadId string) ([]*pb.PollOption, error) {
	var pollOptionsList []models.PollOption
	err := h.DB.WithContext(ctx).Where("thread_id = ?", threadId).Find(&pollOptionsList).Error
	if err != nil {
		return nil, err
	}

	var pollOptionsResponse []*pb.PollOption
	for _, pollOption := range pollOptionsList {
		pollOptionsResponse = append(pollOptionsResponse, &pb.PollOption{
			Option:    pollOption.Option,
			VoteCount: int32(pollOption.VoteCount),
		})
	}

	return pollOptionsResponse, nil
}

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

func (h *Handler) Thread_GetAllThreads(ctx context.Context, req *pb.GetAllThreadsRequest) (*pb.ApiResponseThread, error) {
	zap.L().Info("User getting all threads", zap.String("user_id", req.UserId))

	var threads []models.Thread
	err := h.DB.WithContext(ctx).
		Where("is_scheduled = false OR scheduled_at < ?", time.Now()).
		Order("created_at desc").
		Find(&threads).Error
	if err != nil {
		return &pb.ApiResponseThread{Success: false, Message: "Failed to get threads."}, nil
	}

	if len(threads) == 0 {
		return &pb.ApiResponseThread{Success: false, Message: "No threads found."}, nil
	}

	getAllThreadResponse := &pb.GetThreadsResponse{
		Threads: make([]*pb.Thread, len(threads)),
	}

	for i, request := range threads {
		replyCount := h.getReplyCount(ctx, request.ThreadId)
		likeCount := h.getLikeCount(ctx, request.ThreadId)
		repostCount := h.getRepostCount(ctx, request.ThreadId)

		media, err := h.getMedia(ctx, request.ThreadId)
		if err != nil {
			return &pb.ApiResponseThread{
				Success: false,
				Message: "Failed to fetch media.",
			}, nil
		}

		pollOptions, err := h.getPollOptions(ctx, request.ThreadId)
		if err != nil {
			return &pb.ApiResponseThread{
				Success: false,
				Message: "Failed to fetch poll options.",
			}, nil
		}

		isUserLikingThread := h.IsUserLikingThread(request.ThreadId, req.UserId)
		isUserRepostingThread := h.IsUserRepostingThread(request.ThreadId, req.UserId)
		isUserBookmarkingThread := h.IsUserBookmarkingThread(request.ThreadId, req.UserId)

		getAllThreadResponse.Threads[i] = &pb.Thread{
			ThreadId:        request.ThreadId,
			UserId:          request.UserId,
			Content:         request.Content,
			Category:        request.Category,
			PostedAt:        request.CreatedAt.Format("2006-01-02 15:04:05"),
			ReplyCount:      int32(replyCount),
			LikeCount:       int32(likeCount),
			RepostCount:     int32(repostCount),
			Media:           media,
			PollOptions:     pollOptions,
			IsAdvertisement: request.IsAdvertisement,
			ReplyPermission: request.ReplyPermission,
			Pinned:          request.Pinned,
			IsPrivate:       request.IsPrivate,
			IsLiking:        isUserLikingThread,
			IsReposting:     isUserRepostingThread,
			IsBookmarking:   isUserBookmarkingThread,
		}
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

func (h *Handler) Thread_CreateThread(ctx context.Context, req *pb.PostThread) (*pb.ApiResponseThread, error) {
	zap.L().Info("Creating new thread", zap.String("user_id", req.UserId))

	if req.Title == "" || req.Content == "" || req.Category == "" || req.PollCount < 0 || req.MediaCount < 0 || req.ReplyPermission == "" {
		return &pb.ApiResponseThread{Success: false, Message: "Invalid request parameters."}, nil
	}

	generatedId := uuid.New().String()

	thread := models.Thread{
		ThreadId:        generatedId,
		UserId:          req.UserId,
		CommunityId:     nil,
		Content:         req.Content,
		Category:        req.Category,
		IsPoll:          req.PollCount > 0,
		IsPrivate:       req.IsPrivate,
		IsScheduled:     req.IsScheduled,
		IsAdvertisement: req.IsAdvertisement,
		HasMedia:        len(req.MediaUrls) > 0,
		ReplyPermission: req.ReplyPermission,
		Pinned:          false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if req.IsScheduled && req.ScheduledAt != "" {
		t, err := time.Parse("2006-01-02 15:04:05", req.ScheduledAt)
		if err != nil {
			return &pb.ApiResponseThread{Success: false, Message: "Invalid schedule format"}, nil
		}
		thread.ScheduledAt = &t
	}

	if err := h.DB.WithContext(ctx).Create(&thread).Error; err != nil {
		return &pb.ApiResponseThread{Success: false, Message: "Failed to create thread"}, nil
	}

	for _, url := range req.MediaUrls {
		media := models.Media{
			ThreadId:  generatedId,
			MediaURL:  url,
			MediaType: "image",
		}
		h.DB.WithContext(ctx).Create(&media)
	}

	for _, opt := range req.PollOptions {
		poll := models.PollOption{
			ThreadId:  generatedId,
			Option:    opt,
			VoteCount: 0,
		}
		h.DB.WithContext(ctx).Create(&poll)
	}

	return &pb.ApiResponseThread{
		Success: true,
		Message: "Thread created successfully",
	}, nil
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
		return &pb.ApiResponseThread{Success: true, Message: "Thread liked successfully."}, nil
	}

	return &pb.ApiResponseThread{Success: false, Message: "Unknown error occurred."}, nil
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
