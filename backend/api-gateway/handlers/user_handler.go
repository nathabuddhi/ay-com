package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	pb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/user"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/supabase"
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

	forwardRequest[pb.LoginRequest, pb.String](w, &req, func(ctx context.Context, in *pb.LoginRequest) (*pb.ApiResponse, error) {
		return client.User_Login(ctx, in)
	})
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

	forwardRequest[pb.RegisterRequest, pb.String](w, &req, func(ctx context.Context, in *pb.RegisterRequest) (*pb.ApiResponse, error) {
		resp, err := client.User_Register(ctx, in)
		if err != nil {
			return &pb.ApiResponse{
				Success: false,
				Message: "Error during registration: " + err.Error(),
				Data:    nil,
			}, nil
		}

		if resp.Success {
			var userIdString pb.String
			if err := anypb.UnmarshalTo(resp.Data, &userIdString, proto.UnmarshalOptions{}); err != nil {
				zap.L().Error("Failed to unmarshal response data", zap.Error(err))
				return &pb.ApiResponse{
					Success: false,
					Message: "Failed to process response data",
					Data:    nil,
				}, nil
			}

			go func() {
				if err := supabase.UploadAvatarAndBanner(userIdString.Value, avatarFile, bannerFile); err != nil {
					zap.L().Error("Error uploading avatar and banner to Supabase", zap.Error(err))
				}
			}()
		}

		return resp, nil
	})
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

	var req pb.GetUserRequest
	req.UserId = userID

	forwardRequest[pb.GetUserRequest, pb.UserProfile](w, &req, func(ctx context.Context, in *pb.GetUserRequest) (*pb.ApiResponse, error) {
		resp, err := client.User_GetProfile(ctx, in)
		if err != nil {
			return &pb.ApiResponse{
				Success: false,
				Message: "An error occurred: " + err.Error(),
				Data:    nil,
			}, nil
		}
		return resp, nil
	})
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

	forwardRequest[pb.VerificationRequest, pb.String](w, &req, func(ctx context.Context, in *pb.VerificationRequest) (*pb.ApiResponse, error) {
		return client.RequestVerificationCode(ctx, in)
	})
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

	forwardRequest[pb.ValidateCodeRequest, pb.String](w, &req, func(ctx context.Context, in *pb.ValidateCodeRequest) (*pb.ApiResponse, error) {
		return client.ValidateVerificationCode(ctx, in)
	})
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

	forwardRequest[pb.ChangePasswordRequest, pb.String](w, &req, func(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.ApiResponse, error) {
		return client.ChangePassword(ctx, in)
	})
}
