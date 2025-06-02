package handlers

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-thread/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-thread/proto/thread"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (h *Handler) getVoteCount(ctx context.Context, threadId string, option string) int64 {
	var count int64
	err := h.DB.WithContext(ctx).
		Model(&models.PollVote{}).
		Where("thread_id = ? AND option = ?", threadId, option).
		Count(&count).Error

	if err != nil {
		return 0
	}

	return count
}

func (h *Handler) getPollOptions(ctx context.Context, threadId string, requesterId string) ([]*pb.PollOption, error) {
	var pollOptionsList []models.PollOption
	err := h.DB.WithContext(ctx).Where("thread_id = ?", threadId).Find(&pollOptionsList).Error
	if err != nil {
		return nil, err
	}

	var pollOptionsResponse []*pb.PollOption
	for _, pollOption := range pollOptionsList {
		isUserVoting := false
		var vote models.PollVote
		err := h.DB.WithContext(ctx).Where("thread_id = ? AND option = ? AND user_id = ?", pollOption.ThreadId, pollOption.Option, requesterId).First(&vote).Error
		if err == nil {
			isUserVoting = true
		}

		voteCount := h.getVoteCount(ctx, pollOption.ThreadId, pollOption.Option)

		pollOptionsResponse = append(pollOptionsResponse, &pb.PollOption{
			Option:    pollOption.Option,
			IsVoting:  isUserVoting,
			VoteCount: int32(voteCount),
		})
	}

	return pollOptionsResponse, nil
}

func (h *Handler) Thread_VoteThread(ctx context.Context, req *pb.SubmitVote) (*pb.ApiResponseThread, error) {
	zap.L().Info("Voting on thread", zap.String("thread_id", req.ThreadId), zap.String("user_id", req.UserId))

	if req.ThreadId == "" || req.UserId == "" || req.Content == "" {
		return &pb.ApiResponseThread{

			Success: false,
			Message: "Invalid request parameters.",
		}, nil
	}

	var vote models.PollVote
	err := h.DB.WithContext(ctx).Where("thread_id = ? AND user_id = ?", req.ThreadId, req.UserId).First(&vote).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		// GAK ADA JADI BIKIN BARU
		newVote := models.PollVote{
			ThreadId: req.ThreadId,
			Option:   req.Content,
			UserId:   req.UserId,
		}
		err = h.DB.WithContext(ctx).Create(&newVote).Error
		if err != nil {
			return &pb.ApiResponseThread{
				Success: false,
				Message: err.Error(),
			}, nil
		}
		return &pb.ApiResponseThread{
			Success: true,
			Message: "Voted successfully.",
		}, nil
	} else if err != nil {
		// ERROR
		return &pb.ApiResponseThread{
			Success: false,
			Message: err.Error(),
		}, err
	} else if vote.Option == req.Content {
		// VOTENYA SAMA JADI UNVOTE
		err = h.DB.WithContext(ctx).Where("thread_id = ? AND user_id = ?", req.ThreadId, req.UserId).Delete(&models.PollVote{}).Error
		if err != nil {
			return &pb.ApiResponseThread{
				Success: false,
				Message: err.Error(),
			}, err
		}
		return &pb.ApiResponseThread{
			Success: true,
			Message: "Vote removed successfully.",
		}, nil
	} else {
		vote.Option = req.Content
		err = h.DB.WithContext(ctx).Where("thread_id = ? AND user_id = ?", req.ThreadId, req.UserId).Save(&vote).Error

		if err != nil {
			return &pb.ApiResponseThread{
				Success: false,
				Message: err.Error(),
			}, nil
		}

		return &pb.ApiResponseThread{
			Success: true,
			Message: "Vote updated successfully.",
		}, nil
	}
}
