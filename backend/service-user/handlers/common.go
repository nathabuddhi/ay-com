package handlers

import "github.com/nathabuddhi/ay-com/backend/service-user/models"

func (h *Handlers) IsUserBannedOrDeactivated(user models.User) string {
	if user.IsBanned {
		return "User is banned."
	}
	if user.IsDeactivated {
		return "User is deactivated."
	}
	return ""
}

func (h *Handlers) GetUserById(userId string) (models.User, bool) {
	var user models.User
	err := h.DB.
		Where("user_id = ?", userId).
		First(&user).Error
	if err != nil {
		return user, false
	}
	return user, true
}
func (h *Handlers) GetUserByUsername(username string) (models.User, bool) {
	var user models.User
	err := h.DB.
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return user, false
	}
	return user, true
}
