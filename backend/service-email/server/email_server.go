package server

import (
	"context"
	"fmt"

	"github.com/nathabuddhi/ay-com/backend/service-email/email"
	pb "github.com/nathabuddhi/ay-com/backend/service-email/proto"
)

type EmailServer struct {
	pb.UnimplementedEmailServiceServer
}

func NewEmailServer() *EmailServer {
	return &EmailServer{}
}

func (s *EmailServer) SendVerificationEmail(ctx context.Context, req *pb.SendVerificationEmailRequest) (*pb.SendVerificationEmailResponse, error) {
	subject := "Your Verification Code"
	body := fmt.Sprintf("Your verification code is: <b>%s</b><br><br>This code is only valid for <b>5 minutes</b>.<br><i>You may request another code.</i>", req.VerificationCode)

	err := email.SendEmail(req.ToEmail, subject, body)
	if err != nil {
		return &pb.SendVerificationEmailResponse{
			Success: false,
			Message: "Failed to send email",
		}, err
	}

	return &pb.SendVerificationEmailResponse{
		Success: true,
		Message: "Verification email sent successfully",
	}, nil
}

func (s *EmailServer) SendNotificationEmail(ctx context.Context, req *pb.SendNotificationEmailRequest) (*pb.SendNotificationEmailResponse, error) {
	err := email.SendEmail(req.ToEmail, req.Subject, req.Body)
	if err != nil {
		return &pb.SendNotificationEmailResponse{
			Success: false,
			Message: "Failed to send email",
		}, err
	}

	return &pb.SendNotificationEmailResponse{
		Success: true,
		Message: "Notification email sent successfully",
	}, nil
}
