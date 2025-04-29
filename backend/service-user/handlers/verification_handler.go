package handlers

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/nathabuddhi/ay-com/backend/service-user/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-user/proto/user"
	"github.com/nathabuddhi/ay-com/backend/service-user/rabbitmq"
	"go.uber.org/zap"
)

func (h *Handlers) RequestVerificationCode(ctx context.Context, req *pb.VerificationRequest) (*pb.ApiResponse, error) {
	var user models.User
	err := h.DB.WithContext(ctx).
		Where("email = ?", req.Email).
		First(&user).Error
	if err != nil {
		return &pb.ApiResponse{Success: false, Message: "Email isn't registered."}, nil
	}

	if !user.IsDeactivated {
		return &pb.ApiResponse{Success: false, Message: "Email is already active."}, nil
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	err = h.DB.WithContext(ctx).Exec(`
		INSERT INTO verification_codes (email, code) 
		VALUES (?, ?) 
		ON CONFLICT(email) DO UPDATE SET code = excluded.code
	`, req.Email, code).Error
	if err != nil {
		zap.L().Error("Failed to sign token: " + err.Error())
		return &pb.ApiResponse{Success: false, Message: "An unknown error occured. Please try again."}, nil
	}

	body := fmt.Sprintf("Your verification code is: <b>%s</b><br><br>This code is only valid for <b>5 minutes</b>.<br><i>You may request another code.<br>Ignore this email if this wasn't you.</i>", code)

	rabbitmq.PublishEmail(req.Email, "AY.com Verification Code", body)
	return &pb.ApiResponse{Success: true, Message: "Verification code sent successfully."}, nil
}

func (h *Handlers) ValidateVerificationCode(ctx context.Context, req *pb.ValidateCodeRequest) (*pb.ApiResponse, error) {
	var verificationCode models.VerificationCode

	err := h.DB.WithContext(ctx).
		Where("email = ?", req.Email).
		First(&verificationCode).Error
	if err != nil {
		return &pb.ApiResponse{
			Success: false,
			Message: "Verification code not found or expired.",
		}, nil
	}

	if verificationCode.Code != req.Code {
		return &pb.ApiResponse{
			Success: false,
			Message: "Invalid verification code",
		}, nil
	}

	err = h.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ?", req.Email).
		Update("is_deactivated", false).Error
	if err != nil {
		return &pb.ApiResponse{
			Success: false,
			Message: "Failed to activate user account:" + err.Error(),
		}, nil
	}

	rabbitmq.PublishEmail(req.Email,
		"AY.com Account Activation",
		"Congratulations! Your account has been verified and activated. You can now log in and start using our services.",
	)

	return &pb.ApiResponse{
		Success: true,
		Message: "Account verified successfully.",
	}, nil
}
