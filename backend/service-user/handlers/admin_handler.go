package handlers

import (
	"context"

	"github.com/nathabuddhi/ay-com/backend/service-user/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-user/proto/user"
	rabbitmq "github.com/nathabuddhi/ay-com/backend/service-user/rabbitmqsend"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
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

func (h *Handlers) Admin_GetAllVerifyAccountRequest(ctx context.Context, req *pb.StringUser) (*pb.ApiResponseUser, error) {
	zap.L().Info("Admin getting all verification requests")

	var verifyRequests []models.UserVerificationRequest
	err := h.DB.WithContext(ctx).
		Where("status = ?", "pending").
		Order("status, submitted_at desc").
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
			SelfieUrl:   request.SelfieUrl,
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

func (h *Handlers) Admin_ApprovePremiumRequest(ctx context.Context, req *pb.StringUser) (*pb.ApiResponseUser, error) {
	zap.L().Info("Admin approving premium request with id " + req.Value)

	var request models.UserVerificationRequest
	err := h.DB.WithContext(ctx).
		Where("id = ? AND status = ?", req.Value, "pending").
		First(&request).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Verification request not found."}, nil
	}

	request.Status = "approved"
	err = h.DB.WithContext(ctx).Save(&request).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Failed to approve verification request."}, nil
	}

	var user models.User
	err = h.DB.WithContext(ctx).
		Where("user_id = ?", request.UserId).
		First(&user).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "User not found."}, nil
	}

	user.IsVerified = true
	err = h.DB.WithContext(ctx).Save(&user).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Failed to change user verification status."}, nil
	}

	rabbitmq.PublishDeleteRedis("getprofile/" + user.UserId)
	rabbitmq.PublishEmail(user.Email, "Your verification request has been approved", "Congratulations! Your verification request has been approved. You can now enjoy all the benefits of being a verified user!")
	return &pb.ApiResponseUser{Success: true, Message: "Successfully approved user verification request!"}, nil
}

func (h *Handlers) Admin_RejectPremiumRequest(ctx context.Context, req *pb.RejectPremiumRequest) (*pb.ApiResponseUser, error) {
	zap.L().Info("Admin approving premium request with id " + req.Id + " with reason: " + req.Reason)

	var request models.UserVerificationRequest
	err := h.DB.WithContext(ctx).
		Where("id = ? AND status = ?", req.Id, "pending").
		First(&request).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Verification request not found."}, nil
	}

	request.Status = "rejected"
	err = h.DB.WithContext(ctx).Save(&request).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Failed to reject verification request."}, nil
	}

	var user models.User
	err = h.DB.WithContext(ctx).
		Where("user_id = ?", request.UserId).
		First(&user).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "User not found."}, nil
	}

	rabbitmq.PublishEmail(user.Email, "Your verification request has been rejected", "Dear user, your verification request has been rejected. <br> Reason: "+req.Reason+"<br> If you think this is a mistake, please contact our support team for further assistance.<br> You may reapply for verification after addressing the issue.<br> Thank you for your understanding!")
	return &pb.ApiResponseUser{Success: true, Message: "Successfully rejected user verification request!"}, nil
}
