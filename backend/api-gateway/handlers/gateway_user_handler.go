package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/middleware"
	pb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/user"
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

func processUserResponseWithoutPayload(resp *pb.ApiResponseUser, err error, w http.ResponseWriter) {
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	response := &types.ApiResponse{
		Success: resp.Success,
		Message: resp.Message,
		Payload: nil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func processUserResponseWithPayload[T any](resp *pb.ApiResponseUser, err error, w http.ResponseWriter) {
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	decodedObject := new(T)
	if _, ok := any(decodedObject).(proto.Message); !ok {
		zap.L().Error("Type does not implement proto.Message", zap.String("type", fmt.Sprintf("%T", decodedObject)))
		returnErrorResponse(w, "Failed to decode response data.")
		return
	}

	if err := proto.Unmarshal(resp.Data.GetValue(), any(decodedObject).(proto.Message)); err != nil {
		zap.L().Error("Failed to unmarshal protobuf", zap.Error(err))
		returnErrorResponse(w, "Failed to decode response data.")
		return
	}

	response := &types.ApiResponse{
		Success: resp.Success,
		Message: resp.Message,
		Payload: decodedObject,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func processUserRequest[T any](r *http.Request, w http.ResponseWriter) (resp *T, client pb.UserServiceClient) {
	conn := getUserServiceConn()
	client = pb.NewUserServiceClient(conn)

	var req T
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to encode API request", zap.Error(err))
		returnErrorResponse(w, "Failed to encode API request: "+err.Error())
	}
	return &req, client
}

func IsUserAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIdKey).(string)
		if userID == "" {
			response := types.ApiResponse{
				Success: false,
				Message: "User ID not found in context.",
				Payload: nil,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		conn := getUserServiceConn()
		client := pb.NewUserServiceClient(conn)

		var req pb.StringUser
		req.Value = userID
		res, err := client.Admin_IsUserAdmin(r.Context(), &req)
		if err != nil || res == nil || !res.Value {
			response := types.ApiResponse{
				Success: false,
				Message: "User is not admin.",
				Payload: nil,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(response)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// User_Register godoc
// @Summary Register a new user
// @Description Register a new user with email, name, username, password, gender, date of birth, security question, security answer, avatar, and banner
// @Tags user
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "Email"
// @Param name formData string true "Name"
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Param gender formData string true "Gender"
// @Param date_of_birth formData string true "Date of birth"
// @Param security_question formData string true "Security question"
// @Param security_answer formData string true "Security answer"
// @Param avatar formData file true "Avatar image (PNG)"
// @Param banner formData file true "Banner image (PNG)"
// @Success 200 {object} types.ApiResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/register [post]
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
	req.WantsNewsletter = r.FormValue("wants_newsletter") == "true"

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

// User_GetProfile godoc
// @Summary Get user profile by username
// @Description Get the profile of a user by their username
// @Tags user
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} pb.UserProfile
// @Failure 400 {object} types.ApiResponse
// @Router /user/getprofile/{username} [get]
func User_GetProfile(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Get Profile is called.")

	vars := mux.Vars(r)
	username := vars["username"]
	if username == "" {
		zap.L().Error("Username parameter missing")
		returnErrorResponse(w, "Username is required")
		return
	}

	if !checkRedisData("getprofile/"+username, w) {
		conn := getUserServiceConn()
		client := pb.NewUserServiceClient(conn)

		req := &pb.GetProfileRequest{}

		req.UserId = username
		req.RequesterId = r.Context().Value(middleware.UserIdKey).(string)

		ctx, cancel := createContext()
		defer cancel()

		resp, err := client.User_GetProfile(ctx, req)

		processUserResponseWithPayload[pb.UserProfile](resp, err, w)
	}
}

// User_Login godoc
// @Summary Login a user
// @Description Login a user with email and password
// @Tags user
// @Accept json
// @Produce json
// @Param login body pb.LoginRequest true "Login request"
// @Success 200 {object} pb.LoginResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/login [post]
func User_Login(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Login is called.")

	req, client := processUserRequest[pb.LoginRequest](r, w)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_Login(ctx, req)

	processUserResponseWithPayload[pb.LoginResponse](resp, err, w)
}

// User_RequestVerificationCode godoc
// @Summary Request a verification code
// @Description Request a verification code for user verification
// @Tags user
// @Accept json
// @Produce json
// @Param verification body pb.VerificationRequest true "Verification request"
// @Success 200 {object} types.ApiResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/requestverificationcode [post]
func User_RequestVerificationCode(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Request Verification Code is called.")
	req, client := processUserRequest[pb.VerificationRequest](r, w)
	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.User_RequestVerificationCode(ctx, req)
	processUserResponseWithoutPayload(resp, err, w)
}

// User_ValidateVerificationCode godoc
// @Summary Validate a verification code
// @Description Validate a verification code provided by the user
// @Tags user
// @Accept json
// @Produce json
// @Param validateCode body pb.ValidateCodeRequest true "Validate code request"
// @Success 200 {object} types.ApiResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/validateverificationcode [post]
func User_ValidateVerificationCode(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Validate Verification Code is called.")
	req, client := processUserRequest[pb.ValidateCodeRequest](r, w)
	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.User_ValidateVerificationCode(ctx, req)
	processUserResponseWithoutPayload(resp, err, w)
}

// User_ChangePassword godoc
// @Summary Change user password
// @Description Change the password for the authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Param changePassword body pb.ChangePasswordRequest true "Change password request"
// @Success 200 {object} types.ApiResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/changepassword [patch]
func User_ChangePassword(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Change Password is called.")
	req, client := processUserRequest[pb.ChangePasswordRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.User_ChangePassword(ctx, req)
	processUserResponseWithoutPayload(resp, err, w)
}

// User_GetSecurityQuestion godoc
// @Summary Get security question
// @Description Retrieve the security question for a user
// @Tags user
// @Accept json
// @Produce json
// @Param getSecurityQuestion body pb.GetSecurityQuestionRequest true "Get security question request"
// @Success 200 {object} pb.StringUser
// @Failure 400 {object} types.ApiResponse
// @Router /user/getsecurityquestion [post]
func User_GetSecurityQuestion(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Get Security Question is called.")
	req, client := processUserRequest[pb.GetSecurityQuestionRequest](r, w)
	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.User_GetSecurityQuestion(ctx, req)
	processUserResponseWithPayload[pb.StringUser](resp, err, w)
}

// User_ValidateSecurityAnswer godoc
// @Summary Validate security answer
// @Description Validate the security answer provided by the user
// @Tags user
// @Accept json
// @Produce json
// @Param validateSecurityAnswer body pb.ValidateSecurityAnswerRequest true "Validate security answer request"
// @Success 200 {object} types.ApiResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/validatesecurityanswer [post]
func User_ValidateSecurityAnswer(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Validate Security Answer is called.")
	req, client := processUserRequest[pb.ValidateSecurityAnswerRequest](r, w)
	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.User_ValidateSecurityAnswer(ctx, req)
	processUserResponseWithoutPayload(resp, err, w)
}

// User_ResetPassword godoc
// @Summary Reset user password
// @Description Reset the password for a user
// @Tags user
// @Accept json
// @Produce json
// @Param resetPassword body pb.ResetPasswordRequest true "Reset password request"
// @Success 200 {object} types.ApiResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/resetpassword [put]
func User_ResetPassword(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Reset Password is called.")
	req, client := processUserRequest[pb.ResetPasswordRequest](r, w)
	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.User_ResetPassword(ctx, req)
	processUserResponseWithoutPayload(resp, err, w)
}

// User_FollowUser godoc
// @Summary Follow a user
// @Description Follow another user
// @Tags user
// @Accept json
// @Produce json
// @Param followUser body pb.FollowUserRequest true "Follow user request"
// @Success 200 {object} types.ApiResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/followuser [post]
func User_FollowUser(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Follow User is called.")
	req, client := processUserRequest[pb.FollowUserRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.User_FollowUser(ctx, req)
	processUserResponseWithoutPayload(resp, err, w)
}

// User_UnFollowUser godoc
// @Summary Unfollow a user
// @Description Unfollow another user
// @Tags user
// @Accept json
// @Produce json
// @Param unfollowUser body pb.UnFollowUserRequest true "Unfollow user request"
// @Success 200 {object} types.ApiResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/unfollowuser [post]
func User_UnFollowUser(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Unfollow User is called.")
	req, client := processUserRequest[pb.UnFollowUserRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.User_UnFollowUser(ctx, req)
	processUserResponseWithoutPayload(resp, err, w)
}

// User_BlockUser godoc
// @Summary Block a user
// @Description Block another user
// @Tags user
// @Accept json
// @Produce json
// @Param blockUser body pb.BlockUserRequest true "Block user request"
// @Success 200 {object} types.ApiResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/blockuser [post]
func User_BlockUser(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Block User is called.")
	req, client := processUserRequest[pb.BlockUserRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.User_BlockUser(ctx, req)
	processUserResponseWithoutPayload(resp, err, w)
}

// User_UnBlockUser godoc
// @Summary Unblock a user
// @Description Unblock a previously blocked user
// @Tags user
// @Accept json
// @Produce json
// @Param unblockUser body pb.UnBlockUserRequest true "Unblock user request"
// @Success 200 {object} types.ApiResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/unblockuser [post]
func User_UnBlockUser(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Unblock is called.")
	req, client := processUserRequest[pb.UnBlockUserRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.User_UnBlockUser(ctx, req)
	processUserResponseWithoutPayload(resp, err, w)
}

// User_GetSettings godoc
// @Summary Get user settings
// @Description Retrieve the settings for the authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Param getSettings body pb.GetSettingsRequest true "Get settings request"
// @Success 200 {object} pb.UserSettings
// @Failure 400 {object} types.ApiResponse
// @Router /user/getsettings [post]
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

// User_UpdateSettings godoc
// @Summary Update user settings
// @Description Update the settings for the authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Param updateSettings body pb.UpdateSettingsRequest true "Update settings request"
// @Success 200 {object} types.ApiResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/updatesettings [patch]
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

// User_GetAllFollowers godoc
// @Summary Get all followers
// @Description Get a list of all followers of a user by their ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} pb.AllFollowersResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/getallfollowers/{id} [get]
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

// User_GetAllFollowing godoc
// @Summary Get all following
// @Description Get a list of all users followed by a user by their ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} pb.AllFollowingResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/getallfollowing/{id} [get]
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

// User_SubmitVerifyAccountRequest godoc
// @Summary Submit account verification request
// @Description Submit a request to verify the authenticated user's account
// @Tags user
// @Accept json
// @Produce json
// @Param verifyAccount body pb.SubmitVerifyAccountRequest true "Verify account request"
// @Success 200 {object} types.ApiResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/submitverifyaccountrequest [post]
func User_SubmitVerifyAccountRequest(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Submit Verify Account Request is called.")

	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		zap.L().Error("Failed to parse multipart form", zap.Error(err))
		returnErrorResponse(w, "Failed to parse multipart form")
		return
	}

	file, header, err := r.FormFile("face")
	if err != nil {
		zap.L().Error("Failed to get face file", zap.Error(err))
		returnErrorResponse(w, "Failed to get face file")
		return
	}

	ext := ""
	if dot := strings.LastIndex(header.Filename, "."); dot != -1 {
		ext = header.Filename[dot:]
	}
	imagePath, err := UploadVerifyRequestImage(file, ext)
	if err != nil {
		zap.L().Error("Error uploading verification request image", zap.Error(err))
		returnErrorResponse(w, "Failed to upload verification request image")
		return
	}

	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)
	ctx, cancel := createContext()
	defer cancel()

	var req pb.SubmitVerifyAccountRequest
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	req.IdentityCardNumber = r.FormValue("identity_card_number")
	req.ReasonText = r.FormValue("reason")
	req.SelfieUrl = imagePath
	resp, err := client.User_SubmitVerifyAccountRequest(ctx, &req)
	processUserResponseWithoutPayload(resp, err, w)
}

// User_GetAllVerifyAccountRequest godoc
// @Summary Get all account verification requests
// @Description Retrieve all account verification requests for the authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Param getVerifyRequests body pb.GetAllVerifyAccountRequest true "Get verification requests"
// @Success 200 {object} pb.GetAllVerifyAccountResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/getallverifyaccountrequest [post]
func User_GetAllVerifyAccountRequest(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Fetch Verify Account Request is called.")
	req, client := processUserRequest[pb.GetAllVerifyAccountRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.User_GetAllVerifyAccountRequest(ctx, req)
	processUserResponseWithPayload[pb.GetAllVerifyAccountResponse](resp, err, w)
}

// User_UpdateProfile godoc
// @Summary Update user profile
// @Description Update the profile of the authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Param updateProfile body pb.UpdateUserProfileRequest true "Update profile request"
// @Success 200 {object} types.ApiResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/updateprofile [patch]
func User_UpdateProfile(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Submit Verify Account Request is called.")
	req, client := processUserRequest[pb.UpdateUserProfileRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.User_UpdateProfile(ctx, req)
	processUserResponseWithoutPayload(resp, err, w)
}

// User_DeactivateAccount godoc
// @Summary Deactivate user account
// @Description Deactivate the authenticated user's account
// @Tags user
// @Accept json
// @Produce json
// @Param deactivateAccount body pb.DeactivateAccountRequest true "Deactivate account request"
// @Success 200 {object} types.ApiResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/deactivateaccount [post]
func User_DeactivateAccount(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Deactivate Account is called.")
	req, client := processUserRequest[pb.DeactivateAccountRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.User_DeactivateAccount(ctx, req)
	processUserResponseWithoutPayload(resp, err, w)
}

func User_SearchUser(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Search Users is called.")
	req, client := processUserRequest[pb.StringUser](r, w)
	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.User_SearchUser(ctx, req)
	processUserResponseWithPayload[pb.SearchPeopleResponse](resp, err, w)
}

// User_RefreshToken godoc
// @Summary Refresh access token with refresh token
// @Description Check if the authenticated user's refresh token is valid and returns a new access token.
// @Tags user
// @Accept json
// @Produce json
// @Param checkToken body pb.StringUser true "Token check request"
// @Success 200 {object} pb.LoginResponse
// @Failure 400 {object} pb.LoginResponse
// @Router /user/refreshtoken [post]
func User_RefreshToken(w http.ResponseWriter, r *http.Request) {
	req, client := processUserRequest[pb.StringUser](r, w)

	zap.L().Info("Checking Refresh Token " + req.Value)
	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_RefreshToken(ctx, req)

	processUserResponseWithPayload[pb.LoginResponse](resp, err, w)
}

func User_GetSelfProfile(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Get Self Profile is called.")

	req, client := processUserRequest[pb.StringUser](r, w)
	req.Value = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_GetSelfProfile(ctx, req)

	processUserResponseWithPayload[pb.UserProfile](resp, err, w)
}

// User_GetAllBlocked godoc
// @Summary Get All User Blocked by the User.
// @Description Get All Blocked Users by the User.
// @Tags user
// @Accept json
// @Produce json
// @Param value path string true "value"
// @Success 200 {object} types.ApiResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/getallblocked [post]
func User_GetAllBlocked(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Get All Blocked Request is called.")

	req, client := processUserRequest[pb.StringUser](r, w)
	req.Value = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_GetAllBlocked(ctx, req)

	processUserResponseWithPayload[pb.AllBlockedUserResponse](resp, err, w)
}

// User_GetUserId godoc
// @Summary Get user profile by user id
// @Description Get the profile of a user by their user id
// @Tags user
// @Accept json
// @Produce json
// @Param user_id path string true "user_id"
// @Success 200 {object} pb.UserProfile
// @Failure 400 {object} types.ApiResponse
// @Router /user/getuserid/{user_id} [post]
func User_GetUserId(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Get User Id is called.")

	req, client := processUserRequest[pb.GetProfileRequest](r, w)

	if !checkRedisData("getprofile/"+req.UserId, w) {
		req.RequesterId = r.Context().Value(middleware.UserIdKey).(string)
		ctx, cancel := createContext()
		defer cancel()

		resp, err := client.User_GetUserId(ctx, req)

		processUserResponseWithPayload[pb.UserProfile](resp, err, w)
	}
}

// User_GetFollowRecommendations godoc
// @Summary Get Follow Recommendations Based On Current user
// @Description Get Follow Recommendations Based On Current user
// @Tags user
// @Accept json
// @Produce json
// @Param user_id path string true "StringUser"
// @Success 200 {object} pb.ApiResponse
// @Failure 400 {object} types.ApiResponse
// @Router /user/getfollowrecommendations [get]
func User_GetFollowRecommendations(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User GetFollowRecommendations is called.")

	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	var req pb.StringUser

	req.Value = r.Context().Value(middleware.UserIdKey).(string)

	if !checkRedisData("getfollowrecommendations/"+req.Value, w) {
		ctx, cancel := createContext()
		defer cancel()

		resp, err := client.User_GetFollowRecommendations(ctx, &req)

		processUserResponseWithPayload[pb.GetFollowRecommendationsResponse](resp, err, w)
	}
}

func User_ReportUser(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Report User is called.")

	req, client := processUserRequest[pb.CreateReportRequest](r, w)
	req.ReporterId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.User_SendReport(ctx, req)

	processUserResponseWithPayload[pb.AllBlockedUserResponse](resp, err, w)
}
