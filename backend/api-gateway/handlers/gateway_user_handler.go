package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/middleware"
	pb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/user"
	redis_client "github.com/nathabuddhi/ay-com/backend/api-gateway/redis"
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
	zap.L().Info("User Login is called.")

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

	if !resp.Success {
		returnErrorResponse(w, "Invalid Credentials.")
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
	zap.L().Info("User Register is called.")

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
	zap.L().Info("User Get Profile is called.")

	vars := mux.Vars(r)
	userID := vars["id"]

	redisProfile := redis_client.GetCache("getprofile/" + userID)

	if redisProfile != nil {
		response := types.ApiResponse{
			Success: true,
			Message: "Get User Profile successful.",
			Payload: redisProfile,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		conn := getUserServiceConn()
		client := pb.NewUserServiceClient(conn)

		if userID == "" {
			zap.L().Error("User ID parameter missing")
			returnErrorResponse(w, "User ID is required")
			return
		}

		zap.L().Info("Fetching profile for user ID: " + userID)

		var req pb.GetProfileRequest
		req.UserId = userID
		req.RequesterId = r.Context().Value(middleware.UserIdKey).(string)

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		resp, err := client.User_GetProfile(ctx, &req)
		if err != nil {
			zap.L().Error("Error forwarding request", zap.Error(err))
			returnErrorResponse(w, "Error forwarding request: "+err.Error())
			return
		}

		if !resp.Success {
			returnErrorResponse(w, resp.Message)
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
}

func User_RequestVerificationCode(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Request Verification Code is called.")

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
	zap.L().Info("User Validate Verification Code is called.")

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
	zap.L().Info("User Change Password is called.")

	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	var req pb.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to decode change password request", zap.Error(err))
		returnErrorResponse(w, "Invalid request payload: "+err.Error())
		return
	}

	userID := r.Context().Value(middleware.UserIdKey).(string)
	req.UserId = userID
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
	zap.L().Info("User Get Security Question is called.")

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
	zap.L().Info("User Validate Security Answer is called.")

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
	zap.L().Info("User Reset Password is called.")

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

func User_FollowUser(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Follow User is called.")

	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	var req pb.FollowUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to decode request body", zap.Error(err))
		returnErrorResponse(w, "Invalid request payload: "+err.Error())
		return
	}

	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.User_FollowUser(ctx, &req)
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

func User_UnFollowUser(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Unfollow User is called.")

	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	var req pb.UnFollowUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to decode request body", zap.Error(err))
		returnErrorResponse(w, "Invalid request payload: "+err.Error())
		return
	}

	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.User_UnFollowUser(ctx, &req)
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

func User_BlockUser(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Block User is called.")

	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	var req pb.BlockUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to decode request body", zap.Error(err))
		returnErrorResponse(w, "Invalid request payload: "+err.Error())
		return
	}

	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.User_BlockUser(ctx, &req)
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

func User_UnBlockUser(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Unblock is called.")

	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	var req pb.UnBlockUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to decode request body", zap.Error(err))
		returnErrorResponse(w, "Invalid request payload: "+err.Error())
		return
	}

	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.User_UnBlockUser(ctx, &req)
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

func User_GetSettings(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Get Settings is called.")

	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	var req pb.GetSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to decode request body", zap.Error(err))
		returnErrorResponse(w, "Invalid request payload: "+err.Error())
		return
	}

	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.User_GetSettings(ctx, &req)
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	binaryData := resp.Data.GetValue()
	userSettings, err := decodeResponse[pb.UserSettings](binaryData)
	if err != nil {
		zap.L().Error("Failed to decode user settings.", zap.Error(err))
		returnErrorResponse(w, "Failed to process response data.")
		return
	}

	response := types.ApiResponse{
		Success: resp.Success,
		Message: resp.Message,
		Payload: userSettings,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func User_UpdateSettings(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Update Settings is called.")

	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	var req pb.UpdateSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to decode request body", zap.Error(err))
		returnErrorResponse(w, "Invalid request payload: "+err.Error())
		return
	}

	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.User_UpdateSettings(ctx, &req)
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

func User_GetAllFollowers(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Get All Followers is called.")

	vars := mux.Vars(r)
	userID := vars["id"]

	redisFollowers := redis_client.GetCache("getallfollowers/" + userID)

	if redisFollowers != nil {
		response := types.ApiResponse{
			Success: true,
			Message: "Get All Followers successful.",
			Payload: redisFollowers,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		conn := getUserServiceConn()
		client := pb.NewUserServiceClient(conn)

		if userID == "" {
			zap.L().Error("User ID parameter missing")
			returnErrorResponse(w, "User ID is required")
			return
		}

		zap.L().Info("Fetching followers for user ID: " + userID)

		var req pb.GetAllFollowersRequest
		req.UserId = userID
		req.RequesterId = r.Context().Value(middleware.UserIdKey).(string)

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		resp, err := client.User_GetAllFollowers(ctx, &req)
		if err != nil {
			zap.L().Error("Error forwarding request", zap.Error(err))
			returnErrorResponse(w, "Error forwarding request: "+err.Error())
			return
		}

		if !resp.Success {
			zap.L().Error("No data found in response")
			returnErrorResponse(w, resp.Message)
			return
		}

		binaryData := resp.Data.GetValue()
		allFollowers, err := decodeResponse[pb.AllFollowersResponse](binaryData)
		if err != nil {
			zap.L().Error("Failed to decode response from service.", zap.Error(err))
			returnErrorResponse(w, "Failed to process service response.")
			return
		}

		response := types.ApiResponse{
			Success: resp.Success,
			Message: resp.Message,
			Payload: allFollowers,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func User_GetAllFollowing(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Get All Following is called.")

	vars := mux.Vars(r)
	userID := vars["id"]

	redisFollowing := redis_client.GetCache("getallfollowing/" + userID)

	if redisFollowing != nil {
		response := types.ApiResponse{
			Success: true,
			Message: "Get All Following successful.",
			Payload: redisFollowing,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		conn := getUserServiceConn()
		client := pb.NewUserServiceClient(conn)

		if userID == "" {
			zap.L().Error("User ID parameter missing")
			returnErrorResponse(w, "User ID is required")
			return
		}

		zap.L().Info("Fetching following for user ID: " + userID)

		var req pb.GetAllFollowingRequest
		req.UserId = userID
		req.RequesterId = r.Context().Value(middleware.UserIdKey).(string)

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		resp, err := client.User_GetAllFollowing(ctx, &req)
		if err != nil {
			zap.L().Error("Error forwarding request", zap.Error(err))
			returnErrorResponse(w, "Error forwarding request: "+err.Error())
			return
		}

		if !resp.Success {
			zap.L().Error("No data found in response")
			returnErrorResponse(w, resp.Message)
			return
		}

		binaryData := resp.Data.GetValue()
		allFollowing, err := decodeResponse[pb.AllFollowingResponse](binaryData)
		if err != nil {
			zap.L().Error("Failed to decode response from service.", zap.Error(err))
			returnErrorResponse(w, "Failed to process service response.")
			return
		}

		response := types.ApiResponse{
			Success: resp.Success,
			Message: resp.Message,
			Payload: allFollowing,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func User_SubmitVerifyAccountRequest(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Submit Verify Account Request is called.")

	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	var req pb.SubmitVerifyAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to decode request body", zap.Error(err))
		returnErrorResponse(w, "Invalid request payload: "+err.Error())
		return
	}

	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.User_SubmitVerifyAccountRequest(ctx, &req)
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

func User_GetAllVerifyAccountRequest(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Get All Verify Account Request is called.")

	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	var req pb.GetAllVerifyAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to decode request body", zap.Error(err))
		returnErrorResponse(w, "Invalid request payload: "+err.Error())
		return
	}

	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.User_GetAllVerifyAccountRequest(ctx, &req)
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	binaryData := resp.Data.GetValue()
	allAccountVerificationRequests, err := decodeResponse[pb.GetAllVerifyAccountResponse](binaryData)
	if err != nil {
		zap.L().Error("Failed to decode response from service.", zap.Error(err))
		returnErrorResponse(w, "Failed to process service response.")
		return
	}

	response := types.ApiResponse{
		Success: resp.Success,
		Message: resp.Message,
		Payload: allAccountVerificationRequests,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
