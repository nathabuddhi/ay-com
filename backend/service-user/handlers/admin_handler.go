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

func (h *Handlers) Admin_GetAllUsers(ctx context.Context, req *pb.StringUser) (*pb.ApiResponseUser, error) {
	zap.L().Info("Admin getting all users")

	var users []models.User
	err := h.DB.WithContext(ctx).
		Find(&users).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Failed to get users."}, nil
	}

	if len(users) == 0 {
		return &pb.ApiResponseUser{Success: false, Message: "No users found."}, nil
	}

	getAllUsersResponse := &pb.AdminUserResponse{
		Users: make([]*pb.AdminUserProfile, len(users)),
	}

	for i, user := range users {
		getAllUsersResponse.Users[i] = &pb.AdminUserProfile{
			UserId:     user.UserId,
			Name:       user.Name,
			Username:   user.Username,
			IsVerified: user.IsVerified,
			IsBanned:   user.IsBanned,
		}
	}

	returnData, err := anypb.New(getAllUsersResponse)
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "An error occured: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Get all users successful.",
		Data:    returnData,
	}, nil
}

func (h *Handlers) Admin_ToggleUserBan(ctx context.Context, req *pb.StringUser) (*pb.ApiResponseUser, error) {
	zap.L().Info("Admin toggling ban status for user " + req.Value)

	var user models.User
	err := h.DB.WithContext(ctx).
		Where("user_id = ?", req.Value).
		First(&user).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "User not found."}, nil
	}

	user.IsBanned = !user.IsBanned
	err = h.DB.WithContext(ctx).Save(&user).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Failed to toggle user ban status."}, nil
	}

	rabbitmq.PublishDeleteRedis("getprofile/" + user.UserId)
	return &pb.ApiResponseUser{Success: true, Message: "Successfully toggled user ban status!"}, nil
}

func (h *Handlers) Admin_SendNewsLetter(ctx context.Context, req *pb.SendNewsLetterRequest) (*pb.ApiResponseUser, error) {
	zap.L().Info("Admin sending newsletter with title: " + req.Title)

	if req.Title == "" || req.Content == "" {
		return &pb.ApiResponseUser{Success: false, Message: "Title and content cannot be empty."}, nil
	}

	var users []models.User
	err := h.DB.WithContext(ctx).
		Where("is_deactivated = ? AND is_banned = ?", false, false).
		Find(&users).Error

	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Failed to fetch users to send newsletter."}, nil
	}
	for _, user := range users {
		err = rabbitmq.PublishSendNotification("newsletter", user.UserId, user.Email, "AY.com Newsletter: "+req.Title, req.Content+"<br><br>Thank you for being a part of our community!<br>Best regards,<br>AY.com Team.", "System")
		if err != nil {
			zap.L().Error("Failed to send newsletter to user "+user.UserId, zap.Error(err))
		} else {
			zap.L().Info("Newsletter sent to user " + user.UserId)
		}
	}

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Newsletter sent successfully to all users.",
		Data:    nil,
	}, nil
}

func (h *Handlers) Admin_GetAllReports(ctx context.Context, req *pb.StringUser) (*pb.ApiResponseUser, error) {
	zap.L().Info("Admin getting all user reports")

	var reports []models.UserReport
	err := h.DB.WithContext(ctx).
		Where("status = ?", "pending").
		Order("status desc, submitted_at desc").
		Find(&reports).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Failed to get user reports."}, nil
	}

	if len(reports) == 0 {
		return &pb.ApiResponseUser{Success: false, Message: "No user reports found."}, nil
	}

	getAllReportsResponse := &pb.GetAllReportsResponse{
		Reports: make([]*pb.Report, len(reports)),
	}

	for i, report := range reports {
		getAllReportsResponse.Reports[i] = &pb.Report{
			ReportId:    report.ReportId,
			ReportedId:  report.ReportedId,
			ReporterId:  report.ReporterId,
			Reason:      report.Reason,
			Status:      report.Status,
			SubmittedAt: report.SubmittedAt.String(),
		}
	}

	returnData, err := anypb.New(getAllReportsResponse)
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "An error occurred: " + err.Error(),
			Data:    nil,
		}, nil
	}

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Get all user reports successful.",
		Data:    returnData,
	}, nil
}

func (h *Handlers) Admin_ApproveReport(ctx context.Context, req *pb.StringUser) (*pb.ApiResponseUser, error) {
	zap.L().Info("Admin approving report with ID: " + req.Value)

	var report models.UserReport
	err := h.DB.WithContext(ctx).
		Where("report_id = ? AND status = ?", req.Value, "pending").
		First(&report).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Report not found."}, nil
	}

	var reportedUser models.User
	err = h.DB.WithContext(ctx).
		Where("user_id = ?", report.ReportedId).
		First(&reportedUser).Error

	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Reported user not found."}, nil
	}

	reportedUser.IsBanned = true
	report.Status = "approved"
	err = h.DB.WithContext(ctx).Save(&report).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Failed to approve report."}, nil
	}
	err = h.DB.WithContext(ctx).Save(&reportedUser).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Failed to ban reported user."}, nil
	}

	rabbitmq.PublishDeleteRedis("getprofile/" + reportedUser.UserId)
	rabbitmq.PublishEmail(reportedUser.Email, "Your account has been banned", "Dear user, your account has been banned due to a report against you. If you believe this is a mistake, please contact our support team for further assistance.")

	var reporter models.User
	err = h.DB.WithContext(ctx).
		Where("user_id = ?", report.ReporterId).
		First(&reporter).Error
	if err == nil {
		rabbitmq.PublishSendNotification("system", reporter.UserId, reporter.Email, "Your report has been approved.", "Dear user, your report against "+reportedUser.Username+" has been approved. The user has been banned from the platform. Thank you for helping us maintain a safe community!", "ADMIN")
	}

	return &pb.ApiResponseUser{Success: true, Message: "Successfully approved report and banned user!"}, nil
}

func (h *Handlers) Admin_RejectReport(ctx context.Context, req *pb.StringUser) (*pb.ApiResponseUser, error) {
	zap.L().Info("Admin rejecting report with ID: " + req.Value)

	var report models.UserReport
	err := h.DB.WithContext(ctx).
		Where("report_id = ? AND status = ?", req.Value, "pending").
		First(&report).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Report not found."}, nil
	}

	report.Status = "rejected"
	err = h.DB.WithContext(ctx).Save(&report).Error
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Failed to reject report."}, nil
	}

	return &pb.ApiResponseUser{Success: true, Message: "Successfully rejected report!"}, nil
}
