package handlers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nathabuddhi/ay-com/backend/service-user/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-user/proto/user"
	"go.uber.org/zap"
)

func (h *Handlers) User_CreateReport(ctx context.Context, req *pb.CreateReportRequest) (*pb.ApiResponseUser, error) {
	zap.L().Info("User Create Report is called.")

	if req.ReportedId == "" || req.ReporterId == "" || req.Reason == "" {
		return &pb.ApiResponseUser{Success: false, Message: "Invalid request parameters."}, nil
	}

	generatedId := uuid.New().String()

	report := models.UserReport{
		ReportId:    generatedId,
		ReportedId:  req.ReportedId,
		ReporterId:  req.ReporterId,
		Reason:      req.Reason,
		SubmittedAt: time.Now(),
		Status:      "pending",
	}

	err := h.DB.WithContext(ctx).Create(&report).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Failed to create report."}, nil
	}

	return &pb.ApiResponseUser{Success: true, Message: "Report created successfully!"}, nil
}
