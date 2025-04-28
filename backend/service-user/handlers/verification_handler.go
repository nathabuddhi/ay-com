package handlers

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/nathabuddhi/ay-com/backend/service-user/models"
	emailpb "github.com/nathabuddhi/ay-com/backend/service-user/proto/email"
	redispb "github.com/nathabuddhi/ay-com/backend/service-user/proto/redis"
	userpb "github.com/nathabuddhi/ay-com/backend/service-user/proto/user"
	"go.uber.org/zap"
)

func (h *Handlers) RequestVerificationCode(ctx context.Context, req *userpb.VerificationRequest) (*userpb.ApiResponse, error) {
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	_, err := h.RedisClient.SetKey(ctx, &redispb.SetKeyRequest{
		Key:               fmt.Sprintf("verify:%s", req.Email),
		Value:             code,
		ExpirationSeconds: 300,
	})
	if err != nil {
		return &userpb.ApiResponse{Success: false, Message: "Failed to store verification code"}, err
	}

	_, err = h.EmailClient.SendVerificationEmail(ctx, &emailpb.SendVerificationEmailRequest{
		ToEmail:          req.Email,
		VerificationCode: code,
	})
	if err != nil {
		return &userpb.ApiResponse{Success: false, Message: "Failed to send verification email"}, err
	}

	return &userpb.ApiResponse{Success: true, Message: "Verification code sent successfully"}, nil
}

func (h *Handlers) ValidateVerificationCode(ctx context.Context, req *userpb.ValidateCodeRequest) (*userpb.ApiResponse, error) {
	resp, err := h.RedisClient.GetKey(ctx, &redispb.GetKeyRequest{
		Key: fmt.Sprintf("verify:%s", req.Email),
	})
	if err != nil || !resp.Found {
		return &userpb.ApiResponse{Success: false, Message: "Verification code not found or expired"}, err
	}

	if resp.Value != req.Code {
		return &userpb.ApiResponse{Success: false, Message: "Invalid verification code"}, nil
	}

	err = h.DB.WithContext(ctx).Model(&models.User{}).Where("email = ?", req.Email).Update("is_deactivated", false).Error
	if err != nil {
		return &userpb.ApiResponse{Success: false, Message: "Failed to activate user account"}, err
	}

	h.RedisClient.DeleteKey(ctx, &redispb.DeleteKeyRequest{
		Key: fmt.Sprintf("verify:%s", req.Email),
	})

	_, err = h.EmailClient.SendNotificationEmail(ctx, &emailpb.SendNotificationEmailRequest{
		ToEmail: req.Email,
		Subject: "Your account has been successfully activated!",
		Body:    "Congratulations! Your account has been verified and activated. You can now log in and start using our services.",
	})
	if err != nil {
		zap.L().Error("Failed to send account activation email: " + err.Error())
	}

	return &userpb.ApiResponse{Success: true, Message: "Account verified successfully,"}, nil
}
