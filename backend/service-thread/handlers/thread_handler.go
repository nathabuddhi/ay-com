package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nathabuddhi/ay-com/backend/service-thread/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-thread/proto/thread"
	"github.com/nathabuddhi/ay-com/backend/service-thread/rabbitmq"
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

func (h *Handler) processThreadResponse(ctx context.Context, thread models.Thread, requesterId string) (pb.Thread, error) {
	replyCount := h.getReplyCount(ctx, thread.ThreadId)
	likeCount := h.getLikeCount(ctx, thread.ThreadId)
	repostCount := h.getRepostCount(ctx, thread.ThreadId)

	media, err := h.getMedia(ctx, thread.ThreadId)
	if err != nil {
		return pb.Thread{}, fmt.Errorf("failed to fetch media")
	}

	pollOptions, err := h.getPollOptions(ctx, thread.ThreadId, requesterId)
	if err != nil {
		return pb.Thread{}, fmt.Errorf("failed to fetch poll options")
	}

	isUserLikingThread := h.IsUserLikingThread(thread.ThreadId, requesterId)
	isUserRepostingThread := h.IsUserRepostingThread(thread.ThreadId, requesterId)
	isUserBookmarkingThread := h.IsUserBookmarkingThread(thread.ThreadId, requesterId)

	return pb.Thread{
		ThreadId:        thread.ThreadId,
		UserId:          thread.UserId,
		Content:         thread.Content,
		Category:        thread.Category,
		PostedAt:        thread.CreatedAt.Format("2006-01-02 15:04:05"),
		ReplyCount:      int32(replyCount),
		LikeCount:       int32(likeCount),
		RepostCount:     int32(repostCount),
		Media:           media,
		PollOptions:     pollOptions,
		IsAdvertisement: thread.IsAdvertisement,
		ReplyPermission: thread.ReplyPermission,
		Pinned:          thread.Pinned,
		IsPrivate:       thread.IsPrivate,
		IsLiking:        isUserLikingThread,
		IsReposting:     isUserRepostingThread,
		IsBookmarking:   isUserBookmarkingThread,
		CommunityId: func() string {
			if thread.CommunityId != nil {
				return *thread.CommunityId
			}
			return ""
		}(),
	}, nil
}

func (h *Handler) Thread_GetAllThreads(ctx context.Context, req *pb.GetAllThreadsRequest) (*pb.ApiResponseThread, error) {
	zap.L().Info("User getting all threads", zap.String("user_id", req.UserId))

	var threads []models.Thread
	err := h.DB.WithContext(ctx).
		Where("is_scheduled = false OR scheduled_at < ?", time.Now()).
		Where("category != ?", "Reply").
		Order("created_at desc").
		Find(&threads).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.ApiResponseThread{Success: true, Message: "No threads found."}, nil
		}
		return &pb.ApiResponseThread{Success: false, Message: err.Error()}, nil
	}

	getAllThreadResponse := &pb.GetThreadsResponse{
		Threads: make([]*pb.Thread, len(threads)),
	}

	for i, request := range threads {
		threadResponse, err := h.processThreadResponse(ctx, request, req.UserId)
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

func (h *Handler) Thread_CreateThread(ctx context.Context, req *pb.PostThread) (*pb.ApiResponseThread, error) {
	zap.L().Info("Creating new thread", zap.String("user_id", req.UserId))

	if req.Content == "" || req.PollCount < 0 || req.MediaCount < 0 || req.ReplyPermission == "" {
		return &pb.ApiResponseThread{Success: false, Message: "Invalid request parameters."}, nil
	}

	generatedId := uuid.New().String()

	thread := models.Thread{
		ThreadId:        generatedId,
		UserId:          req.UserId,
		CommunityId:     &req.CommunityId,
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
			ThreadId: generatedId,
			Option:   opt,
		}
		h.DB.WithContext(ctx).Create(&poll)
	}

	if req.ReplyTo != "" {
		err := h.Thread_ReplyThread(ctx, generatedId, req.ReplyTo)

		if err != nil {
			return &pb.ApiResponseThread{Success: false, Message: "Failed to reply to thread: " + err.Error()}, nil
		}
	}

	re := regexp.MustCompile(`@(\S+)`)
	matches := re.FindAllStringSubmatch(req.Content, -1)
	for _, match := range matches {
		username := strings.TrimSpace(match[1])
		if username != "" {
			rabbitmq.PublishSendMention(username, req.UserId)
		}
	}

	return &pb.ApiResponseThread{
		Success: true,
		Message: "Thread created successfully!",
	}, nil
}

func (h *Handler) Thread_GetThreadById(ctx context.Context, req *pb.GeneralThreadRequest) (*pb.ApiResponseThread, error) {
	zap.L().Info("Getting thread by ID", zap.String("thread_id", req.ThreadId))

	var thread models.Thread
	err := h.DB.WithContext(ctx).Where("thread_id = ?", req.ThreadId).First(&thread).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.ApiResponseThread{Success: false, Message: "Thread not found."}, nil
		}
		return &pb.ApiResponseThread{Success: false, Message: "An error occurred: " + err.Error()}, nil
	}

	threadResponse, err := h.processThreadResponse(ctx, thread, req.ThreadId)
	if err != nil {
		return &pb.ApiResponseThread{
			Success: false,
			Message: "An error occurred: " + err.Error(),
		}, nil
	}

	var comments []models.ThreadReply
	err = h.DB.WithContext(ctx).Where("reply_to_id = ?", req.ThreadId).Find(&comments).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return &pb.ApiResponseThread{Success: false, Message: err.Error()}, nil
	}

	getThreadDetailResponse := &pb.GetThreadDetailResponse{
		Thread:  &threadResponse,
		Replies: make([]*pb.Thread, len(comments)),
	}

	for i, request := range comments {
		var comment models.Thread
		err := h.DB.WithContext(ctx).Where("thread_id = ?", request.ThreadId).First(&comment).Error
		if err == nil {
			commentResponse, err := h.processThreadResponse(ctx, comment, req.UserId)
			if err == nil {
				getThreadDetailResponse.Replies[i] = &commentResponse
			}
		}
		if err != nil {
			return &pb.ApiResponseThread{
				Success: false,
				Message: "Failed to fetch replies: " + err.Error(),
			}, nil
		}
	}

	redisData, err := json.Marshal(getThreadDetailResponse)
	if err == nil {
		rabbitmq.PublishSetRedis("getthread/"+thread.ThreadId, string(redisData))
	}

	returnData, err := anypb.New(getThreadDetailResponse)
	if err != nil {
		return &pb.ApiResponseThread{
			Success: false,
			Message: "An error occurred: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponseThread{
		Success: true,
		Message: "Get thread by ID successful.",
		Data:    returnData,
	}, nil
}

func (h *Handler) Thread_TogglePinThread(ctx context.Context, req *pb.GeneralThreadRequest) (*pb.ApiResponseThread, error) {
	zap.L().Info("User "+req.UserId+" is toggling pin for thread", zap.String("thread_id", req.ThreadId))

	if req.ThreadId == "" || req.UserId == "" {
		return &pb.ApiResponseThread{Success: false, Message: "Invalid request parameters."}, nil
	}

	var thread models.Thread
	err := h.DB.WithContext(ctx).Where("thread_id = ? AND user_id = ?", req.ThreadId, req.UserId).First(&thread).Error
	if err == nil {
		err := h.DB.WithContext(ctx).Model(&thread).Update("pinned", !thread.Pinned).Error
		if err != nil {
			return &pb.ApiResponseThread{Success: false, Message: err.Error()}, nil
		}
		rabbitmq.PublishDeleteRedis("getthread/" + thread.ThreadId)
		return &pb.ApiResponseThread{Success: true, Message: "Pin status toggled successfully."}, nil
	} else if err == gorm.ErrRecordNotFound {
		return &pb.ApiResponseThread{Success: false, Message: "Thread not found."}, nil
	} else {
		return &pb.ApiResponseThread{Success: false, Message: err.Error()}, nil
	}
}

func (h *Handler) Admin_DeleteThread(ctx context.Context, req *pb.StringThread) (*pb.ApiResponseThread, error) {
	zap.L().Info("Admin is deleting thread", zap.String("thread_id", req.Value))

	if req.Value == "" {
		return &pb.ApiResponseThread{Success: false, Message: "Invalid request parameters."}, nil
	}

	var thread models.Thread
	err := h.DB.WithContext(ctx).Where("thread_id = ?", req.Value).First(&thread).Error
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
			Message: "Thread not found.",
		}, nil
	}
}

func (h *Handler) Thread_DeleteThread(ctx context.Context, req *pb.GeneralThreadRequest) (*pb.ApiResponseThread, error) {
	zap.L().Info("User "+req.UserId+" is deleting thread", zap.String("thread_id", req.ThreadId))

	if req.ThreadId == "" || req.UserId == "" {
		return &pb.ApiResponseThread{Success: false, Message: "Invalid request parameters."}, nil
	}

	var thread models.Thread
	err := h.DB.WithContext(ctx).Where("thread_id = ?", req.ThreadId).First(&thread).Error
	if err != nil {
		return &pb.ApiResponseThread{
			Success: false,
			Message: "Thread not found or you do not have permission to delete this thread.",
		}, nil
	}

	if thread.UserId == req.UserId {
		err := h.DB.WithContext(ctx).Delete(&thread).Error
		if err != nil {
			return &pb.ApiResponseThread{Success: false, Message: "Failed to delete thread."}, nil
		}
		rabbitmq.PublishDeleteRedis("getthread/" + thread.ThreadId)
		return &pb.ApiResponseThread{Success: true, Message: "Thread deleted successfully."}, nil
	}

	if thread.Category == "Reply" {
		var reply models.ThreadReply
		err := h.DB.WithContext(ctx).Where("thread_id = ?", thread.ThreadId).First(&reply).Error
		if err == nil && reply.ReplyToId != "" {
			var parentThread models.Thread
			err := h.DB.WithContext(ctx).Where("thread_id = ?", reply.ReplyToId).First(&parentThread).Error
			if err == nil && parentThread.UserId == req.UserId {
				err := h.DB.WithContext(ctx).Delete(&thread).Error
				if err != nil {
					return &pb.ApiResponseThread{Success: false, Message: "Failed to delete thread."}, nil
				}
				rabbitmq.PublishDeleteRedis("getthread/" + thread.ThreadId)
				rabbitmq.PublishDeleteRedis("getthread/" + reply.ReplyToId)
				return &pb.ApiResponseThread{Success: true, Message: "Thread deleted successfully."}, nil
			}
		}
	}

	return &pb.ApiResponseThread{
		Success: false,
		Message: "Thread not found or you do not have permission to delete this thread.",
	}, nil
}

func (h *Handler) Thread_GetUserThreads(ctx context.Context, req *pb.UserToUserRequest) (*pb.ApiResponseThread, error) {
	zap.L().Info("Getting user threads ", zap.String("user_id", req.UserId))

	var threads []models.Thread
	err := h.DB.WithContext(ctx).
		Where("is_scheduled = false OR scheduled_at < ?", time.Now()).
		Where("user_id = ? AND category != ?", req.UserId, "Reply").
		Order("pinned desc, created_at desc").
		Find(&threads).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.ApiResponseThread{Success: true, Message: "No threads found."}, nil
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

func (h *Handler) Thread_GetUserMediaThreads(ctx context.Context, req *pb.UserToUserRequest) (*pb.ApiResponseThread, error) {
	var threads []models.Thread
	err := h.DB.WithContext(ctx).
		Where("is_scheduled = false OR scheduled_at < ?", time.Now()).
		Where("user_id = ? AND has_media = ?", req.UserId, true).
		Order("created_at desc").
		Find(&threads).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.ApiResponseThread{Success: true, Message: "No threads found."}, nil
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

func (h *Handler) Thread_GetCommunityMediaThreads(ctx context.Context, req *pb.GeneralThreadRequest) (*pb.ApiResponseThread, error) {
	var threads []models.Thread
	err := h.DB.WithContext(ctx).
		Where("is_scheduled = false OR scheduled_at < ?", time.Now()).
		Where("community_id = ? AND has_media = ?", req.ThreadId, true).
		Order("created_at desc").
		Find(&threads).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.ApiResponseThread{Success: true, Message: "No threads found."}, nil
		}
		return &pb.ApiResponseThread{Success: false, Message: err.Error()}, nil
	}

	getAllThreadResponse := &pb.GetThreadsResponse{
		Threads: make([]*pb.Thread, len(threads)),
	}

	for i, request := range threads {
		threadResponse, err := h.processThreadResponse(ctx, request, req.UserId)
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

func (h *Handler) Thread_GetTrendingHashtags(ctx context.Context, req *pb.StringThread) (*pb.ApiResponseThread, error) {
	zap.L().Info("Getting trending hashtags")

	var threads []models.Thread
	err := h.DB.WithContext(ctx).
		Where("is_scheduled = false OR scheduled_at < ?", time.Now()).
		Find(&threads).Error
	if err != nil {
		return &pb.ApiResponseThread{Success: false, Message: err.Error()}, nil
	}

	hashtagCount := make(map[string]int)
	re := regexp.MustCompile(`#(\w+)`)

	for _, thread := range threads {
		matches := re.FindAllStringSubmatch(thread.Content, -1)
		for _, match := range matches {
			if len(match) > 1 {
				hashtag := strings.ToLower(match[1])
				hashtagCount[hashtag]++
			}
		}
	}

	type hashtagStat struct {
		Tag   string
		Count int
	}
	var stats []hashtagStat
	for tag, count := range hashtagCount {
		stats = append(stats, hashtagStat{Tag: tag, Count: count})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Count > stats[j].Count
	})

	topN := 10
	if len(stats) < topN {
		topN = len(stats)
	}

	var trending []*pb.Hashtag
	for _, stat := range stats[:topN] {
		trending = append(trending, &pb.Hashtag{
			Hashtag:     stat.Tag,
			ThreadCount: int32(stat.Count),
		})
	}

	resp := &pb.GetTrendingHashtagsResponse{
		Hashtags: trending,
	}

	data, err := anypb.New(resp)
	if err != nil {
		return &pb.ApiResponseThread{Success: false, Message: err.Error()}, nil
	}

	return &pb.ApiResponseThread{
		Success: true,
		Message: "Top 10 trending hashtags fetched successfully.",
		Data:    data,
	}, nil
}

func (h *Handler) Thread_GetFollowingThreads(ctx context.Context, req *pb.GetFollowingThreadRequest) (*pb.ApiResponseThread, error) {
	zap.L().Info("User getting following threads", zap.String("user_id", req.UserId))

	if len(req.FollowedIds) == 0 {
		return &pb.ApiResponseThread{Success: true, Message: "No followed users."}, nil
	}

	var userThreads []models.Thread
	err := h.DB.WithContext(ctx).
		Where("user_id IN ?", req.FollowedIds).
		Where("is_scheduled = false OR scheduled_at < ?", time.Now()).
		Where("category != ?", "Reply").
		Find(&userThreads).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return &pb.ApiResponseThread{Success: false, Message: err.Error()}, nil
	}

	var reposts []models.ThreadRepost
	err = h.DB.WithContext(ctx).
		Where("user_id IN ?", req.FollowedIds).
		Find(&reposts).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return &pb.ApiResponseThread{Success: false, Message: err.Error()}, nil
	}

	repostedThreadIDs := make([]string, 0, len(reposts))
	repostMap := make(map[string]models.ThreadRepost)
	for _, repost := range reposts {
		repostedThreadIDs = append(repostedThreadIDs, repost.ThreadId)
		repostMap[repost.ThreadId] = repost
	}

	var repostThreads []models.Thread
	if len(repostedThreadIDs) > 0 {
		err = h.DB.WithContext(ctx).
			Where("thread_id IN ? AND user_id != ?", repostedThreadIDs, req.UserId).
			Find(&repostThreads).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return &pb.ApiResponseThread{Success: false, Message: err.Error()}, nil
		}
	}

	allThreads := append(userThreads, repostThreads...)

	sort.Slice(allThreads, func(i, j int) bool {
		return allThreads[i].CreatedAt.After(allThreads[j].CreatedAt)
	})

	getThreadsResponse := &pb.GetThreadsResponse{
		Threads: make([]*pb.Thread, 0, len(allThreads)),
	}

	for _, thread := range allThreads {
		threadResp, err := h.processThreadResponse(ctx, thread, req.UserId)
		if err != nil {
			return &pb.ApiResponseThread{
				Success: false,
				Message: "Failed to process thread: " + err.Error(),
			}, nil
		}

		getThreadsResponse.Threads = append(getThreadsResponse.Threads, &threadResp)
	}

	returnData, err := anypb.New(getThreadsResponse)
	if err != nil {
		return &pb.ApiResponseThread{
			Success: false,
			Message: "An error occurred: " + err.Error(),
		}, nil
	}

	return &pb.ApiResponseThread{
		Success: true,
		Message: "Get following threads successful.",
		Data:    returnData,
	}, nil
}

func (h *Handler) Thread_GetAdvertisementThreads(ctx context.Context, req *pb.GeneralThreadRequest) (*pb.ApiResponseThread, error) {
	zap.L().Info("Getting advertisement threads")

	var threads []models.Thread
	err := h.DB.WithContext(ctx).
		Where("is_scheduled = false OR scheduled_at < ?", time.Now()).
		Where("is_advertisement = ?", true).
		Order("created_at desc").
		Find(&threads).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.ApiResponseThread{Success: true, Message: "No advertisement threads found."}, nil
		}
		return &pb.ApiResponseThread{Success: false, Message: err.Error()}, nil
	}

	getAllThreadResponse := &pb.GetThreadsResponse{
		Threads: make([]*pb.Thread, len(threads)),
	}

	for i, request := range threads {
		threadResponse, err := h.processThreadResponse(ctx, request, req.UserId)
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
			Message: "An error occurred: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponseThread{
		Success: true,
		Message: "Get advertisement threads successful.",
		Data:    returnData,
	}, nil
}

func (h *Handler) Thread_GetCommunityThreads(ctx context.Context, req *pb.GeneralThreadRequest) (*pb.ApiResponseThread, error) {
	zap.L().Info("Getting advertisement threads")

	var threads []models.Thread
	err := h.DB.WithContext(ctx).
		Where("is_scheduled = false OR scheduled_at < ?", time.Now()).
		Where("community_id = ?", req.ThreadId).
		Order("created_at desc").
		Find(&threads).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.ApiResponseThread{Success: true, Message: "No advertisement threads found."}, nil
		}
		return &pb.ApiResponseThread{Success: false, Message: err.Error()}, nil
	}

	getAllThreadResponse := &pb.GetThreadsResponse{
		Threads: make([]*pb.Thread, len(threads)),
	}

	for i, request := range threads {
		threadResponse, err := h.processThreadResponse(ctx, request, req.UserId)
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
			Message: "An error occurred: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponseThread{
		Success: true,
		Message: "Get advertisement threads successful.",
		Data:    returnData,
	}, nil
}

func (h *Handler) Thread_GetThreadByhashtag(ctx context.Context, req *pb.GeneralThreadRequest) (*pb.ApiResponseThread, error) {
	zap.L().Info("Getting threads by hashtag", zap.String("hashtag", req.ThreadId))

	if req.ThreadId == "" {
		return &pb.ApiResponseThread{Success: false, Message: "Invalid request parameters."}, nil
	}

	var threads []models.Thread
	err := h.DB.WithContext(ctx).
		Where("is_scheduled = false OR scheduled_at < ?", time.Now()).
		Where("content LIKE ?", "%#"+req.ThreadId+"%").
		Order("created_at desc").
		Find(&threads).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.ApiResponseThread{Success: true, Message: "No threads found."}, nil
		}
		return &pb.ApiResponseThread{Success: false, Message: err.Error()}, nil
	}

	getAllThreadResponse := &pb.GetThreadsResponse{
		Threads: make([]*pb.Thread, len(threads)),
	}

	for i, request := range threads {
		threadResponse, err := h.processThreadResponse(ctx, request, req.UserId)
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
			Message: "An error occurred: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponseThread{
		Success: true,
		Message: "Get threads by hashtag successful.",
		Data:    returnData,
	}, nil
}
