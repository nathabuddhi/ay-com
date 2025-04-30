package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	pb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/user"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/supabase"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/types"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	userServiceConn *grpc.ClientConn
	userServiceOnce sync.Once
)

func getUserServiceConn() *grpc.ClientConn {
	userServiceOnce.Do(func() {
		var err error
		userServiceConn, err = grpc.NewClient(USER_SERVICE_PATH, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			zap.L().Fatal("Failed to connect to user service: " + err.Error())
		}
	})
	return userServiceConn
}

func User_Login(w http.ResponseWriter, r *http.Request) {
	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	var req pb.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to decode login request", zap.Error(err))
		returnErrorResponse(w, "Invalid request payload: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.User_Login(ctx, &req)
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	if resp.Data == nil {
		zap.L().Error("No data found in response")
		returnErrorResponse(w, "No profile data found")
		return
	}

	binaryData := resp.Data.GetValue()

	jwtToken, err := decodeResponse[pb.UserProfile](binaryData)
	if err != nil {
		zap.L().Error("Failed to decode JWT Token.", zap.Error(err))
		returnErrorResponse(w, "Failed to process JWT token.")
		return
	}

	response := types.ApiResponse{
		Success: resp.Success,
		Message: resp.Message,
		Payload: jwtToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func User_Register(w http.ResponseWriter, r *http.Request) {
	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		zap.L().Error("Failed to parse multipart form", zap.Error(err))
		returnErrorResponse(w, "Failed to parse multipart form")
		return
	}

	avatarFile, avatarHeader, err := r.FormFile("avatar")
	if err != nil {
		zap.L().Error("Failed to get avatar file", zap.Error(err))
		returnErrorResponse(w, "Failed to get avatar file")
		return
	}
	defer avatarFile.Close()

	bannerFile, bannerHeader, err := r.FormFile("banner")
	if err != nil {
		zap.L().Error("Failed to get banner file", zap.Error(err))
		returnErrorResponse(w, "Failed to get banner file")
		return
	}
	defer bannerFile.Close()

	if len(avatarHeader.Filename) < 4 || strings.ToLower(avatarHeader.Filename[len(avatarHeader.Filename)-4:]) != ".png" {
		returnErrorResponse(w, "Avatar file must be a .png file")
		return
	}

	if len(bannerHeader.Filename) < 4 || strings.ToLower(bannerHeader.Filename[len(bannerHeader.Filename)-4:]) != ".png" {
		returnErrorResponse(w, "Banner file must be a .png file")
		return
	}

	var req pb.RegisterRequest
	req.Email = r.FormValue("email")
	req.Name = r.FormValue("name")
	req.Username = r.FormValue("username")
	req.Password = r.FormValue("password")
	req.Gender = r.FormValue("gender")
	req.DateOfBirth = r.FormValue("date_of_birth")
	req.SecurityQuestion = r.FormValue("security_question")
	req.SecurityAnswer = r.FormValue("security_answer")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.User_Register(ctx, &req)
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	if resp.Success {
		var userIdString pb.String
		if err := anypb.UnmarshalTo(resp.Data, &userIdString, proto.UnmarshalOptions{}); err != nil {
			zap.L().Error("Failed to unmarshal response data", zap.Error(err))
			returnErrorResponse(w, "Failed to process response data")
			return
		}

		go func() {
			if err := supabase.UploadAvatarAndBanner(userIdString.Value, avatarFile, bannerFile); err != nil {
				zap.L().Error("Error uploading avatar and banner to Supabase", zap.Error(err))
			}
		}()
	}
	response := types.ApiResponse{
		Success: resp.Success,
		Message: resp.Message,
		Payload: nil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func User_GetProfile(w http.ResponseWriter, r *http.Request) {
	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	vars := mux.Vars(r)
	userID := vars["id"]

	if userID == "" {
		zap.L().Error("User ID parameter missing")
		returnErrorResponse(w, "User ID is required")
		return
	}

	zap.L().Info("Fetching profile for user ID: " + userID)

	var req pb.GetProfileRequest
	req.UserId = userID

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.User_GetProfile(ctx, &req)
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	if resp.Data == nil {
		zap.L().Error("No data found in response")
		returnErrorResponse(w, "No profile data found")
		return
	}

	binaryData := resp.Data.GetValue()
	userProfile, err := decodeResponse[pb.UserProfile](binaryData)
	if err != nil {
		zap.L().Error("Failed to decode user profile", zap.Error(err))
		returnErrorResponse(w, "Failed to process user profile data.")
		return
	}

	response := types.ApiResponse{
		Success: resp.Success,
		Message: resp.Message,
		Payload: userProfile,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func User_RequestVerificationCode(w http.ResponseWriter, r *http.Request) {
	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	var req pb.VerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to decode verification request", zap.Error(err))
		returnErrorResponse(w, "Invalid request payload: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.User_RequestVerificationCode(ctx, &req)
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	response := types.ApiResponse{
		Success: resp.Success,
		Message: resp.Message,
		Payload: nil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func User_ValidateVerificationCode(w http.ResponseWriter, r *http.Request) {
	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	var req pb.ValidateCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to decode validate code request", zap.Error(err))
		returnErrorResponse(w, "Invalid request payload: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.User_ValidateVerificationCode(ctx, &req)
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	response := types.ApiResponse{
		Success: resp.Success,
		Message: resp.Message,
		Payload: nil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func User_ChangePassword(w http.ResponseWriter, r *http.Request) {
	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	var req pb.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to decode change password request", zap.Error(err))
		returnErrorResponse(w, "Invalid request payload: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.User_ChangePassword(ctx, &req)
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	response := types.ApiResponse{
		Success: resp.Success,
		Message: resp.Message,
		Payload: resp.Data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func User_GetSecurityQuestion(w http.ResponseWriter, r *http.Request) {
	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	var req pb.GetSecurityQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to decode change password request", zap.Error(err))
		returnErrorResponse(w, "Invalid request payload: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.User_GetSecurityQuestion(ctx, &req)
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	binaryData := resp.Data.GetValue()
	securityQuestion, err := decodeResponse[pb.String](binaryData)
	if err != nil {
		zap.L().Error("Failed to decode security question", zap.Error(err))
		returnErrorResponse(w, "Failed to process security question data.")
		return
	}

	response := types.ApiResponse{
		Success: resp.Success,
		Message: resp.Message,
		Payload: securityQuestion,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func User_ValidateSecurityAnswer(w http.ResponseWriter, r *http.Request) {
	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	var req pb.ValidateSecurityAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to decode change password request", zap.Error(err))
		returnErrorResponse(w, "Invalid request payload: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.User_ValidateSecurityAnswer(ctx, &req)
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	response := types.ApiResponse{
		Success: resp.Success,
		Message: resp.Message,
		Payload: nil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func User_ResetPassword(w http.ResponseWriter, r *http.Request) {
	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	var req pb.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to decode change password request", zap.Error(err))
		returnErrorResponse(w, "Invalid request payload: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.User_ResetPassword(ctx, &req)
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	response := types.ApiResponse{
		Success: resp.Success,
		Message: resp.Message,
		Payload: nil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
