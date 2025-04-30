package handlers

import "github.com/nathabuddhi/ay-com/backend/service-user/models"

func (h *Handlers) GetFollowers(user_id string) (int, error) {
	var count int64
	err := h.DB.Model(&models.UserFollower{}).Where("user_id = ?", user_id).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (h *Handlers) GetFollowing(user_id string) (int, error) {
	var count int64
	err := h.DB.Model(&models.UserFollower{}).Where("follower_id = ?", user_id).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
