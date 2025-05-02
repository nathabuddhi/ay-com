package handlers

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/middleware"
	pb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/user"
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

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_Register(ctx, &req)
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	if resp.Success {
		var userIdString pb.StringUser
		if err := anypb.UnmarshalTo(resp.Data, &userIdString, proto.UnmarshalOptions{}); err != nil {
			zap.L().Error("Failed to unmarshal response data", zap.Error(err))
			returnErrorResponse(w, "Failed to process response data")
			return
		}

		if err := UploadAvatar(userIdString.Value, avatarFile); err != nil {
			zap.L().Error("Error uploading avatar.", zap.Error(err))
			returnErrorResponse(w, "Failed to upload avatar files.")
			return
		}

		if err := UploadBanner(userIdString.Value, bannerFile); err != nil {
			zap.L().Error("Error uploading banner.", zap.Error(err))
			returnErrorResponse(w, "Failed to upload banner files.")
			return
		}
	}
	processUserResponseWithoutPayload(resp, err, w)
}

func User_GetProfile(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Get Profile is called.")

	vars := mux.Vars(r)
	userID := vars["id"]
	if userID == "" {
		zap.L().Error("User ID parameter missing")
		returnErrorResponse(w, "User ID is required")
		return
	}

	if !checkRedisData("getprofile/"+userID, w) {
		req, client := processUserRequest[pb.GetProfileRequest](r, w)
		req.UserId = userID
		req.RequesterId = r.Context().Value(middleware.UserIdKey).(string)

		ctx, cancel := createContext()
		defer cancel()

		resp, err := client.User_GetProfile(ctx, req)

		processUserResponseWithPayload[pb.UserProfile](resp, err, w)
	}
}

func User_Login(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Login is called.")

	req, client := processUserRequest[pb.LoginRequest](r, w)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_Login(ctx, req)

	processUserResponseWithPayload[pb.StringUser](resp, err, w)
}

func User_RequestVerificationCode(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Request Verification Code is called.")

	req, client := processUserRequest[pb.VerificationRequest](r, w)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_RequestVerificationCode(ctx, req)

	processUserResponseWithoutPayload(resp, err, w)
}

func User_ValidateVerificationCode(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Validate Verification Code is called.")

	req, client := processUserRequest[pb.ValidateCodeRequest](r, w)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_ValidateVerificationCode(ctx, req)

	processUserResponseWithoutPayload(resp, err, w)
}

func User_ChangePassword(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Change Password is called.")

	req, client := processUserRequest[pb.ChangePasswordRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_ChangePassword(ctx, req)

	processUserResponseWithoutPayload(resp, err, w)
}

func User_GetSecurityQuestion(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Get Security Question is called.")

	req, client := processUserRequest[pb.GetSecurityQuestionRequest](r, w)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_GetSecurityQuestion(ctx, req)

	processUserResponseWithPayload[pb.StringUser](resp, err, w)
}

func User_ValidateSecurityAnswer(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Validate Security Answer is called.")

	req, client := processUserRequest[pb.ValidateSecurityAnswerRequest](r, w)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_ValidateSecurityAnswer(ctx, req)

	processUserResponseWithoutPayload(resp, err, w)
}

func User_ResetPassword(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Reset Password is called.")

	req, client := processUserRequest[pb.ResetPasswordRequest](r, w)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_ResetPassword(ctx, req)

	processUserResponseWithoutPayload(resp, err, w)
}

func User_FollowUser(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Follow User is called.")

	req, client := processUserRequest[pb.FollowUserRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_FollowUser(ctx, req)

	processUserResponseWithoutPayload(resp, err, w)
}

func User_UnFollowUser(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Unfollow User is called.")

	req, client := processUserRequest[pb.UnFollowUserRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_UnFollowUser(ctx, req)

	processUserResponseWithoutPayload(resp, err, w)
}

func User_BlockUser(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Block User is called.")

	req, client := processUserRequest[pb.BlockUserRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_BlockUser(ctx, req)

	processUserResponseWithoutPayload(resp, err, w)
}

func User_UnBlockUser(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Unblock is called.")

	req, client := processUserRequest[pb.UnBlockUserRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_UnBlockUser(ctx, req)

	processUserResponseWithoutPayload(resp, err, w)
}

func User_GetSettings(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Get Settings is called.")

	req, client := processUserRequest[pb.GetSettingsRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_GetSettings(ctx, req)
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	processUserResponseWithPayload[pb.UserSettings](resp, err, w)
}

func User_UpdateSettings(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Update Settings is called.")

	req, client := processUserRequest[pb.UpdateSettingsRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_UpdateSettings(ctx, req)
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	processUserResponseWithoutPayload(resp, err, w)
}

func User_GetAllFollowers(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Get All Followers is called.")

	vars := mux.Vars(r)
	userID := vars["id"]

	if !checkRedisData("getallfollowers/"+userID, w) {
		req, client := processUserRequest[pb.GetAllFollowersRequest](r, w)
		req.RequesterId = r.Context().Value(middleware.UserIdKey).(string)
		req.UserId = userID

		ctx, cancel := createContext()
		defer cancel()

		resp, err := client.User_GetAllFollowers(ctx, req)

		processUserResponseWithPayload[pb.AllFollowersResponse](resp, err, w)
	}
}

func User_GetAllFollowing(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Get All Following is called.")

	vars := mux.Vars(r)
	userID := vars["id"]

	if !checkRedisData("getallfollowing/"+userID, w) {
		req, client := processUserRequest[pb.GetAllFollowingRequest](r, w)
		req.RequesterId = r.Context().Value(middleware.UserIdKey).(string)
		req.UserId = userID

		ctx, cancel := createContext()
		defer cancel()

		resp, err := client.User_GetAllFollowing(ctx, req)

		processUserResponseWithPayload[pb.AllFollowingResponse](resp, err, w)
	}
}

func User_SubmitVerifyAccountRequest(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Submit Verify Account Request is called.")

	req, client := processUserRequest[pb.SubmitVerifyAccountRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_SubmitVerifyAccountRequest(ctx, req)

	processUserResponseWithoutPayload(resp, err, w)
}

func User_GetAllVerifyAccountRequest(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Submit Verify Account Request is called.")

	req, client := processUserRequest[pb.GetAllVerifyAccountRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_GetAllVerifyAccountRequest(ctx, req)

	processUserResponseWithPayload[pb.GetAllVerifyAccountResponse](resp, err, w)
}

func User_UpdateProfile(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Submit Verify Account Request is called.")

	req, client := processUserRequest[pb.UpdateUserProfileRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_UpdateProfile(ctx, req)

	processUserResponseWithoutPayload(resp, err, w)
}

func User_IsUserPrivate(user_id string) (bool, error) {
	zap.L().Info("User Is User Private is called.")

	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)
	req := &pb.IsAccountPrivateRequest{}
	req.UserId = user_id

	ctx, cancel := createContext()
	defer cancel()

	isPrivate, err := client.User_IsUserPrivate(ctx, req)

	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		return false, err
	}
	return isPrivate.Value, nil
}
