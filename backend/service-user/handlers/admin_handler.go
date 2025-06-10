package handlers

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-user/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-user/proto/user"
	"go.uber.org/zap"
)

func (h *Handlers) Admin_IsUserAdmin(ctx context.Context, req *pb.StringUser) (*pb.BoolUser, error) {
	zap.L().Info("Checking if user " + req.Value + " is admin.")
	var user models.User
	err := h.DB.WithContext(ctx).
		Where("user_id = ? AND is_admin = ?", req.Value, true).
		First(&user).Error
	if err != nil {
		zap.L().Error("Error checking if user is admin", zap.Error(err))
		return &pb.BoolUser{Value: false}, nil
	}
	return &pb.BoolUser{Value: true}, nil
}
