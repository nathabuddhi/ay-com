package handlers

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nathabuddhi/ay-com/backend/service-user/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-user/proto/user"
	rabbitmq "github.com/nathabuddhi/ay-com/backend/service-user/rabbitmqsend"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/anypb"
)

func (h *Handlers) User_RequestVerificationCode(ctx context.Context, req *pb.VerificationRequest) (*pb.ApiResponseUser, error) {
	zap.L().Info("User requesting verification code", zap.String("email", req.Email))

	var user models.User
	err := h.DB.WithContext(ctx).
		Where("email = ?", req.Email).
		First(&user).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Email isn't registered."}, nil
	}

	if !user.IsDeactivated {
		return &pb.ApiResponseUser{Success: false, Message: "Email is already active."}, nil
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	err = h.DB.WithContext(ctx).Exec(`
		INSERT INTO verification_codes (email, code, expiry) 
		VALUES (?, ?, ?) 
		ON CONFLICT(email) DO UPDATE SET code = excluded.code, expiry = excluded.expiry
	`, req.Email, code, time.Now().Add(time.Minute*15)).Error
	if err != nil {
		zap.L().Error("Failed to sign token: " + err.Error())
		return &pb.ApiResponseUser{Success: false, Message: "An unknown error occured. Please try again."}, nil
	}

	body := fmt.Sprintf("Your verification code is: <b>%s</b><br><br>This code is only valid for <b>15 minutes</b>.<br><i>You may request another code.<br>Ignore this email if this wasn't you.</i>", code)

	rabbitmq.PublishEmail(req.Email, "AY.com Verification Code", body)
	return &pb.ApiResponseUser{Success: true, Message: "Verification code sent successfully."}, nil
}

func (h *Handlers) User_ValidateVerificationCode(ctx context.Context, req *pb.ValidateCodeRequest) (*pb.ApiResponseUser, error) {
	zap.L().Info("User validating verification code", zap.String("email", req.Email))

	var user models.User
	err := h.DB.WithContext(ctx).
		Where("email = ?", req.Email).
		First(&user).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Email isn't registered."}, nil
	}

	if !user.IsDeactivated {
		return &pb.ApiResponseUser{Success: false, Message: "Email is already active."}, nil
	}

	var verificationCode models.VerificationCode
	err = h.DB.WithContext(ctx).
		Where("email = ?", req.Email).
		First(&verificationCode).Error
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "It doesn't seem you have requested a verification code.",
		}, nil
	}

	if verificationCode.Code != req.Code {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "Invalid verification code.",
		}, nil
	}

	if time.Now().After(verificationCode.Expiry) {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "Verification code expired.",
		}, nil
	}

	err = h.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ?", req.Email).
		Update("is_deactivated", false).Error
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "Failed to activate user account:" + err.Error(),
		}, nil
	}

	h.DB.Delete(&models.VerificationCode{}, "email = ?", req.Email)

	rabbitmq.PublishSendNotification("system", user.UserId, req.Email,
		"AY.com Account Activation",
		"Congratulations! Your account has been verified and activated. You can now log in and start using our services.",
		"System")

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Account verified successfully.",
	}, nil
}

func (h *Handlers) User_SubmitVerifyAccountRequest(ctx context.Context, req *pb.SubmitVerifyAccountRequest) (*pb.ApiResponseUser, error) {
	zap.L().Info("User submitting verification request", zap.String("user_id", req.UserId))

	var user models.User
	err := h.DB.WithContext(ctx).
		Where("user_id = ?", req.UserId).
		First(&user).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Email isn't registered."}, nil
	}

	if user.IsVerified {
		return &pb.ApiResponseUser{Success: false, Message: "You are already verified."}, nil
	}

	if user.IsDeactivated {
		return &pb.ApiResponseUser{Success: false, Message: "Account is not active."}, nil
	}

	if user.IsBanned {
		return &pb.ApiResponseUser{Success: false, Message: "Account is not active."}, nil
	}

	var verifyRequest models.UserVerificationRequest
	err = h.DB.WithContext(ctx).
		Where("user_id = ? AND status IN ('pending', 'accepted')", req.UserId).
		First(&verifyRequest).Error

	if err == nil {
		return &pb.ApiResponseUser{Success: false, Message: "You already have a pending or accepted verification request."}, nil
	}

	hashedIdCardNumber, err := bcrypt.GenerateFromPassword([]byte(req.IdentityCardNumber), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Error("Failed to hash identity card number", zap.Error(err))
		return &pb.ApiResponseUser{Success: false, Message: "Failed to process identity card number."}, nil
	}

	var existingRequests []models.UserVerificationRequest
	err = h.DB.WithContext(ctx).
		Where("status IN ('pending', 'accepted')").
		Find(&existingRequests).Error
	if err != nil {
		zap.L().Error("Failed to check existing id card numbers", zap.Error(err))
		return &pb.ApiResponseUser{Success: false, Message: "Failed to check identity card number usage."}, nil
	}

	for _, reqItem := range existingRequests {
		if bcrypt.CompareHashAndPassword([]byte(reqItem.IdentityCardNumber), []byte(req.IdentityCardNumber)) == nil {
			return &pb.ApiResponseUser{Success: false, Message: "This identity card number has already been used for a verification request. If this isn't you, contact support and local authority immediately."}, nil
		}
	}

	var generatedId string
	for {
		generatedId = uuid.New().String()
		if err := h.DB.Where("id = ?", generatedId).First(&verifyRequest).Error; err != nil {
			break
		}
	}

	verifyRequest.Id = generatedId
	verifyRequest.UserId = req.UserId
	verifyRequest.IdentityCardNumber = string(hashedIdCardNumber)
	verifyRequest.SelfieUrl = req.SelfieUrl
	verifyRequest.ReasonText = req.ReasonText
	verifyRequest.Status = "pending"
	verifyRequest.SubmittedAt = time.Now()

	err = h.DB.WithContext(ctx).Create(&verifyRequest).Error
	if err != nil {
		zap.L().Error("Failed to create verification request", zap.Error(err))
		return &pb.ApiResponseUser{Success: false, Message: "Failed to create verification request: " + err.Error()}, nil
	}

	idCardLength := len(req.IdentityCardNumber)
	firstPart := req.IdentityCardNumber[:idCardLength/4]
	lastPart := req.IdentityCardNumber[(idCardLength*3)/4:]
	maskedIdCard := firstPart + strings.Repeat("*", idCardLength-len(firstPart)-len(lastPart)) + lastPart

	rabbitmq.PublishEmail(user.Email, "AY.com Account Premium Verification Request", "Your verification request has been submitted successfully. We will review it and get back to you shortly. <br><br>Request Id: "+verifyRequest.Id+"<br>Identity Card Number: "+maskedIdCard+"<br>Submitted At: "+verifyRequest.SubmittedAt.String()+"<br><br>If this wasn't you contact support immediately.")

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Account Verification request submitted successfully.",
		Data:    nil,
	}, nil
}

func (h *Handlers) User_GetAllVerifyAccountRequest(ctx context.Context, req *pb.GetAllVerifyAccountRequest) (*pb.ApiResponseUser, error) {
	zap.L().Info("User getting all verification requests", zap.String("user_id", req.UserId))

	var user models.User
	err := h.DB.WithContext(ctx).
		Where("user_id = ?", req.UserId).
		First(&user).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Invalid Credentials."}, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "Invalid Credentials.",
		}, nil
	}

	if user.IsVerified {
		return &pb.ApiResponseUser{Success: false, Message: "Email is already verified."}, nil
	}

	if user.IsDeactivated {
		return &pb.ApiResponseUser{Success: false, Message: "Account is not active."}, nil
	}

	if user.IsBanned {
		return &pb.ApiResponseUser{Success: false, Message: "Account is not active."}, nil
	}

	var verifyRequests []models.UserVerificationRequest
	err = h.DB.WithContext(ctx).
		Where("user_id = ?", req.UserId).
		Order("submitted_at desc").
		Find(&verifyRequests).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Failed to get verification requests."}, nil
	}

	if len(verifyRequests) == 0 {
		return &pb.ApiResponseUser{Success: false, Message: "No verification requests found."}, nil
	}

	getAllVerifyAccountResponse := &pb.GetAllVerifyAccountResponse{
		Requests: make([]*pb.VerifyAccountRequest, len(verifyRequests)),
	}

	for i, request := range verifyRequests {
		getAllVerifyAccountResponse.Requests[i] = &pb.VerifyAccountRequest{
			Id:          request.Id,
			UserId:      request.UserId,
			Status:      request.Status,
			SubmittedAt: request.SubmittedAt.String(),
			ReasonText:  request.ReasonText,
		}
	}

	returnData, err := anypb.New(getAllVerifyAccountResponse)
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "An error occured: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Get verification requests successful.",
		Data:    returnData,
	}, nil
}
