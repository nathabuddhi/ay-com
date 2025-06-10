package handlers

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-thread/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-thread/proto/thread"
	"go.uber.org/zap"
)

func (h *Handler) Thread_GetThreadCategories(ctx context.Context, req *pb.StringThread) (*pb.ThreadCategories, error) {
	zap.L().Info("Getting thread categories for thread ID: " + req.Value)

	var categories []models.ThreadCategory
	err := h.DB.WithContext(ctx).
		Find(&categories).Error
	if err != nil {
		zap.L().Error("Error fetching thread categories", zap.Error(err))
		return &pb.ThreadCategories{
			Categories: []string{"General", "Gaming", "Technology", "Lifestyle", "Education"},
		}, nil
	}

	threadCategories := &pb.ThreadCategories{
		Categories: make([]string, len(categories)),
	}

	for i, category := range categories {
		threadCategories.Categories[i] = category.Category
	}

	return threadCategories, nil
}

func (h *Handler) Admin_AddCategory(ctx context.Context, req *pb.StringThread) (*pb.ApiResponseThread, error) {
	zap.L().Info("Adding new thread category", zap.String("category", req.Value))

	if err := h.DB.WithContext(ctx).Create(
		&models.ThreadCategory{
			Category: req.Value,
		}).Error; err != nil {
		zap.L().Error("Error adding thread category", zap.Error(err))
		return &pb.ApiResponseThread{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.ApiResponseThread{
		Success: true,
		Message: "Category added successfully.",
	}, nil
}

func (h *Handler) Admin_DeleteCategory(ctx context.Context, req *pb.StringThread) (*pb.ApiResponseThread, error) {
	zap.L().Info("Deleting thread category", zap.String("category", req.Value))

	if err := h.DB.WithContext(ctx).Where("category = ?", req.Value).Delete(&models.ThreadCategory{}).Error; err != nil {
		zap.L().Error("Error deleting thread category", zap.Error(err))
		return &pb.ApiResponseThread{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.ApiResponseThread{
		Success: true,
		Message: "Category deleted successfully.",
	}, nil
}
