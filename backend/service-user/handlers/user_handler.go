package handlers

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-user/models"
	"gorm.io/gorm"
)

type Handlers struct {
	DB *gorm.DB
}

func NewHandlers(db *gorm.DB) *Handlers {
	return &Handlers{
		DB: db,
	}
}

func (h *Handlers) IsUserFollowing(ctx context.Context, userId string, targetUserId string) bool {
	var count int64
	err := h.DB.WithContext(ctx).
		Model(&models.UserFollowing{}).
		Where("user_id = ? AND followed_id = ?", userId, targetUserId).
		Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}
