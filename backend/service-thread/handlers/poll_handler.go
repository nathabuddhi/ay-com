package handlers

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-thread/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-thread/proto/thread"
)

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
