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

func (h *Handler) Thread_GetAllThreads(ctx context.Context, req *pb.GetAllThreadsRequest) (*pb.ApiResponseThread, error) {
	zap.L().Info("User getting all verification requests", zap.String("user_id", req.UserId))

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
		// Get dynamic counts for reply, like, repost
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
