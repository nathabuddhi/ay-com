package handlers

import (
	"context"
	"fmt"

	"github.com/nathabuddhi/ay-com/backend/service-thread/models"
	"github.com/nathabuddhi/ay-com/backend/service-thread/rabbitmq"
	"go.uber.org/zap"
)

func (h *Handler) getReplyCount(ctx context.Context, threadId string) int {
	var count int64
	h.DB.WithContext(ctx).Model(&models.ThreadReply{}).Where("thread_id = ?", threadId).Count(&count)
	return int(count)
}

func (h *Handler) Thread_ReplyThread(ctx context.Context, threadId string, replyToId string) error {
	zap.L().Info("Thread " + threadId + " is replying to thread " + replyToId)
	if threadId == "" || replyToId == "" {
		return fmt.Errorf("no thread or reply provided")
	}

	reply := &models.ThreadReply{
		ThreadId:  threadId,
		ReplyToId: replyToId,
	}

	if err := h.DB.WithContext(ctx).Create(reply).Error; err != nil {
		return err
	}

	rabbitmq.PublishDeleteRedis("getthread/" + threadId)

	return nil
}
