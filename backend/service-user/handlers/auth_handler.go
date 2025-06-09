package handlers

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nathabuddhi/ay-com/backend/service-user/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-user/proto/user"
	rabbitmq "github.com/nathabuddhi/ay-com/backend/service-user/rabbitmqsend"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/anypb"
)

func (h *Handlers) CreateAccessToken(userId string) (string, error) {
	var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
	})
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

func (h *Handlers) CreateRefreshToken(userId string) (string, error) {
	var jwtKey = []byte(os.Getenv("JWT_REFRESH_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		zap.L().Error("Failed to sign token: " + err.Error())
		return "", err
	}

	err = h.DB.Exec(`
		INSERT INTO refresh_tokens (user_id, token) 
		VALUES (?, ?) 
		ON CONFLICT(user_id) DO UPDATE SET token = excluded.token
	`, userId, tokenString).Error

	return tokenString, err
}

func (h *Handlers) User_RefreshToken(ctx context.Context, req *pb.StringUser) (*pb.ApiResponseUser, error) {
	var jwtKey = []byte(os.Getenv("JWT_REFRESH_KEY"))

	token, err := jwt.Parse(req.Value, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "Invalid Token.",
			Data:    nil,
		}, nil
	}

	userID := claims["user_id"].(string)
	var count int64
	err = h.DB.Table("refresh_tokens").
		Where("user_id = ? AND token = ?", userID, req.Value).
		Count(&count).Error
	if err != nil || count != 1 || !token.Valid {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "Invalid Token.",
			Data:    nil,
		}, nil
	}

	var user models.User
	if err := h.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "Invalid token.",
		}, nil
	}

	newAccessToken, err := h.CreateAccessToken(user.UserId)
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "An Error Occured: " + err.Error(),
		}, nil
	}

	newRefreshToken, err := h.CreateRefreshToken(user.UserId)
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "An Error Occured: " + err.Error(),
		}, nil
	}

	refreshTokenResponse := &pb.LoginResponse{
		Token:        newAccessToken,
		RefreshToken: newRefreshToken,
		UserId:       user.UserId,
		Username:     user.Username,
		Name:         user.Name,
		IsVerified:   user.IsVerified,
		IsAdmin:      user.IsAdmin,
	}
	dataReturn, err := anypb.New(refreshTokenResponse)
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "An Error Occured: " + err.Error(),
		}, nil
	}

	zap.L().Info("User " + user.UserId + " refreshed token.")

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Token refreshed successfully.",
		Data:    dataReturn,
	}, nil
}

func (h *Handlers) User_Login(ctx context.Context, req *pb.LoginRequest) (*pb.ApiResponseUser, error) {
	var user models.User

	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "Invalid Credentials.",
		}, nil
	}

	if user.IsDeactivated {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "Account is not active. Please verify your email or contact support.",
		}, nil
	}

	if user.IsBanned {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "Your account is banned. Contact support if you think this is a mistake.",
		}, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "Invalid Credentials.",
		}, nil
	}

	accessToken, err := h.CreateAccessToken(user.UserId)
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "An Error Occured: " + err.Error(),
		}, nil
	}

	refreshToken, err := h.CreateRefreshToken(user.UserId)
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "An Error Occured: " + err.Error(),
		}, nil
	}

	loginResponse := &pb.LoginResponse{
		UserId:       user.UserId,
		Username:     user.Username,
		Name:         user.Name,
		IsVerified:   user.IsVerified,
		Token:        accessToken,
		RefreshToken: refreshToken,
		IsAdmin:      user.IsAdmin,
	}
	dataReturn, err := anypb.New(loginResponse)
	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "An Error Occured: " + err.Error(),
		}, nil
	}

	zap.L().Info("User " + user.UserId + " Logged in.")

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Login successful.",
		Data:    dataReturn,
	}, nil
}

func (h *Handlers) User_Register(ctx context.Context, req *pb.RegisterRequest) (*pb.ApiResponseUser, error) {
	if req.Email == "" || req.Username == "" || req.Password == "" || req.Name == "" || req.SecurityQuestion == "" || req.SecurityAnswer == "" || req.Gender == "" || req.DateOfBirth == "" {
		return &pb.ApiResponseUser{Success: false, Message: "All fields must be filled."}, nil
	}

	if len(req.Name) < 5 || !regexp.MustCompile(`^[A-Za-z\s]+$`).MatchString(req.Name) {
		return &pb.ApiResponseUser{Success: false, Message: "Name must be more than 4 characters and contain only letters and spaces."}, nil
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.com$`).MatchString(req.Email) {
		return &pb.ApiResponseUser{Success: false, Message: "Invalid email format. Must end with .com"}, nil
	}

	var existingUser models.User

	if err := h.DB.
		Where("username = ? OR email = ?", req.Username, req.Email).
		First(&existingUser).Error; err == nil {

		if existingUser.Username == req.Username && !existingUser.IsDeactivated {
			return &pb.ApiResponseUser{Success: false, Message: "Username is already taken."}, nil
		}

		if existingUser.Email == req.Email {
			if existingUser.IsDeactivated {
				var verifCode models.VerificationCode
				if err := h.DB.
					Where("email = ? AND expiry > ?", existingUser.Email, time.Now().Add(-time.Hour).Format("2006-01-02 15:04:05")).
					First(&verifCode).Error; err == nil {
					return &pb.ApiResponseUser{Success: false, Message: "Email is already registered."}, nil
				} else {
					h.DB.Delete(&existingUser)
				}
			} else {
				return &pb.ApiResponseUser{Success: false, Message: "Email is already registered."}, nil
			}
		}
	}

	if len(req.Password) < 8 ||
		!regexp.MustCompile(`[A-Z]`).MatchString(req.Password) ||
		!regexp.MustCompile(`[a-z]`).MatchString(req.Password) ||
		!regexp.MustCompile(`[0-9]`).MatchString(req.Password) ||
		!regexp.MustCompile(`[!@#~$%^&*()+|_]`).MatchString(req.Password) {
		return &pb.ApiResponseUser{Success: false, Message: "Password must be at least 8 characters long and include uppercase, lowercase, number, and special character."}, nil
	}

	if req.Gender != "male" && req.Gender != "female" {
		return &pb.ApiResponseUser{Success: false, Message: "Gender must be 'male' or 'female'."}, nil
	}

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Invalid date of birth format. Use YYYY-MM-DD."}, nil
	}
	if time.Since(dob).Hours() < 13*365*24 {
		return &pb.ApiResponseUser{Success: false, Message: "You must be at least 13 years old to register."}, nil
	}

	validQuestions := map[string]bool{
		"What was the name of your first pet?":    true,
		"What city were you born in?":             true,
		"What is your favorite video game?":       true,
		"What was the name of your first school?": true,
		"What was your childhood nickname?":       true,
	}
	if !validQuestions[req.SecurityQuestion] {
		return &pb.ApiResponseUser{Success: false, Message: "Invalid security question selected."}, nil
	}

	// if err := verifyRecaptcha(req.RecaptchaToken); err != nil {
	// 	return &pb.ApiResponseUser{Success: false, Message: "reCAPTCHA verification failed."}, nil
	// }

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var generatedId string
	for {
		generatedId = uuid.New().String()
		if err := h.DB.Where("user_id = ?", generatedId).First(&existingUser).Error; err != nil {
			break
		}
	}

	user := models.User{
		UserId:           generatedId,
		Name:             req.Name,
		Username:         req.Username,
		Email:            req.Email,
		Password:         string(passwordHash),
		Gender:           req.Gender,
		DateOfBirth:      dob,
		SecurityQuestion: req.SecurityQuestion,
		SecurityAnswer:   req.SecurityAnswer,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		JoinedAt:         time.Now(),
		IsVerified:       false,
		IsDeactivated:    true,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Failed to register user: " + err.Error()}, nil
	}

	h.User_RequestVerificationCode(ctx, &pb.VerificationRequest{
		Email: req.Email,
	})

	anyGeneratedId, _ := anypb.New(&pb.StringUser{Value: generatedId})
	return &pb.ApiResponseUser{Success: true, Message: "User registered successfully. Please verify your email.", Data: anyGeneratedId}, nil
}

func (h *Handlers) User_ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ApiResponseUser, error) {
	if req.Email == "" || req.OldPassword == "" || req.NewPassword == "" || req.UserId == "" {
		return &pb.ApiResponseUser{Success: false, Message: "All fields must be filled."}, nil
	}

	if len(req.NewPassword) < 8 ||
		!regexp.MustCompile(`[A-Z]`).MatchString(req.NewPassword) ||
		!regexp.MustCompile(`[a-z]`).MatchString(req.NewPassword) ||
		!regexp.MustCompile(`[0-9]`).MatchString(req.NewPassword) ||
		!regexp.MustCompile(`[!@#~$%^&*()+|_]`).MatchString(req.NewPassword) {
		return &pb.ApiResponseUser{Success: false, Message: "Password must be at least 8 characters long and include uppercase, lowercase, number, and special character."}, nil
	}

	var user models.User
	if err := h.DB.Where("email = ? AND user_id = ?", req.Email, req.UserId).First(&user).Error; err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "User not found."}, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.NewPassword)); err == nil {
		return &pb.ApiResponseUser{Success: false, Message: "Old and new password may not be the same!"}, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Old password is incorrect."}, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Error hashing new password."}, err
	}

	if err := h.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Failed to update password."}, err
	}

	rabbitmq.PublishEmail(req.Email,
		"AY.com Password Change.",
		"Your Password was just changed at "+time.Now().String()+".<br>If this wasn't you, contact support immediately.",
	)

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Password changed successfully.",
	}, nil
}

func (h *Handlers) User_GetSecurityQuestion(ctx context.Context, req *pb.GetSecurityQuestionRequest) (*pb.ApiResponseUser, error) {
	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "User not found."}, nil
	}

	if user.IsDeactivated {
		return &pb.ApiResponseUser{Success: false, Message: "Account is not active. Please verify your email or contact support."}, nil
	}

	if user.IsBanned {
		return &pb.ApiResponseUser{Success: false, Message: "Your account is banned. Contact support if you think this is a mistake."}, nil
	}

	anyQuestion, err := anypb.New(&pb.StringUser{Value: user.SecurityQuestion})
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Failed to get security question."}, err
	}

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Security question retrieved successfully.",
		Data:    anyQuestion,
	}, nil
}

func (h *Handlers) User_ValidateSecurityAnswer(ctx context.Context, req *pb.ValidateSecurityAnswerRequest) (*pb.ApiResponseUser, error) {
	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "User not found."}, nil
	}

	if user.SecurityAnswer != req.Answer {
		return &pb.ApiResponseUser{Success: false, Message: "Invalid Credentials."}, nil
	}

	code := fmt.Sprintf("%10d", rand.Intn(10000000000))

	err := h.DB.WithContext(ctx).Exec(`
		INSERT INTO reset_password_codes (email, code, expiry) 
		VALUES (?, ?, ?) 
		ON CONFLICT(email) DO UPDATE SET code = excluded.code, expiry = excluded.expiry
	`, req.Email, code, time.Now().Add(time.Minute*5)).Error

	if err != nil {
		return &pb.ApiResponseUser{
			Success: false,
			Message: "An error occured. Please try again or contact support.\n" + err.Error(),
			Data:    nil,
		}, nil
	}

	link := "http://localhost:5173/reset-password"
	rabbitmq.PublishEmail(req.Email,
		"AY.com Password Reset.",
		"Please go to <a href='"+link+"' target='_blank'>reset-password-page</a> to reset your password and insert the following code: <b>"+code+"</b>.<br> <i>This code is valid for 10 minutes.</i>",
	)

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Security answer verified. Check your email for further instructions.",
		Data:    nil,
	}, nil
}

func (h *Handlers) User_ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ApiResponseUser, error) {
	if req.Email == "" || req.Code == "" || req.NewPassword == "" {
		return &pb.ApiResponseUser{Success: false, Message: "All fields must be filled."}, nil
	}

	if len(req.NewPassword) < 8 ||
		!regexp.MustCompile(`[A-Z]`).MatchString(req.NewPassword) ||
		!regexp.MustCompile(`[a-z]`).MatchString(req.NewPassword) ||
		!regexp.MustCompile(`[0-9]`).MatchString(req.NewPassword) ||
		!regexp.MustCompile(`[!@#~$%^&*()+|_]`).MatchString(req.NewPassword) {
		return &pb.ApiResponseUser{Success: false, Message: "Password must be at least 8 characters long and include uppercase, lowercase, number, and special character."}, nil
	}

	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "User not found."}, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.NewPassword)); err == nil {
		return &pb.ApiResponseUser{Success: false, Message: "Old and new password may not be the same!"}, nil
	}

	var count int

	err := h.DB.WithContext(ctx).Raw(`
		SELECT COUNT(*) AS count
		FROM reset_password_codes
		WHERE email = ? AND code = ? AND expiry > ? 
	`, req.Email, req.Code, time.Now().Format("2006-01-02 15:04:05")).Scan(&count).Error

	if err != nil || count != 1 {
		return &pb.ApiResponseUser{Success: false, Message: "Invalid or expired code."}, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Error hashing new password."}, err
	}

	if err := h.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		return &pb.ApiResponseUser{Success: false, Message: "Failed to update password."}, err
	}

	h.DB.WithContext(ctx).Exec(`
		DELETE FROM reset_password_codes
		WHERE email = ?
	`, req.Email)

	rabbitmq.PublishEmail(req.Email,
		"AY.com Password Change.",
		"Your Password was just changed at "+time.Now().String()+".<br>If this wasn't you, contact support immediately.",
	)

	return &pb.ApiResponseUser{
		Success: true,
		Message: "Password reset successfully.",
		Data:    nil,
	}, nil
}

func (h *Handlers) User_CheckToken(ctx context.Context, req *pb.StringUser) (*pb.BoolUser, error) {
	var user models.User
	if err := h.DB.Where("user_id = ? AND is_deactivated = false AND is_banned = FALSE", req.Value).First(&user).Error; err != nil {
		zap.L().Info("Token " + req.Value + " is invalid.")
		return &pb.BoolUser{Value: false}, nil
	} else {
		zap.L().Info("Token " + req.Value + " is valid.")
		return &pb.BoolUser{Value: true}, nil
	}
}
