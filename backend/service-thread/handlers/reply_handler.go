package handlers

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-thread/models"
)

func (h *Handler) getReplyCount(ctx context.Context, threadId string) int {
	var count int64
	h.DB.WithContext(ctx).Model(&models.ThreadReply{}).Where("thread_id = ?", threadId).Count(&count)
	return int(count)
}
