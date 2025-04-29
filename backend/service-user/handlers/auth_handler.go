package handlers

import (
	"context"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nathabuddhi/ay-com/backend/service-user/models"
	pb "github.com/nathabuddhi/ay-com/backend/service-user/proto/user"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/anypb"
)

func (h *Handlers) Login(ctx context.Context, req *pb.LoginRequest) (*pb.ApiResponse, error) {
	var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

	var user models.User

	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return &pb.ApiResponse{
			Success: false,
			Message: "Invalid Credentials.",
		}, nil
	}

	if user.IsDeactivated {
		return &pb.ApiResponse{
			Success: false,
			Message: "Account is not active. Please verify your email or contact support.",
		}, nil
	}

	if user.IsBanned {
		return &pb.ApiResponse{
			Success: false,
			Message: "Your account is banned. Contact support if you think this is a mistake.",
		}, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return &pb.ApiResponse{
			Success: false,
			Message: "Invalid Credentials.",
		}, nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserId,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return &pb.ApiResponse{
			Success: false,
			Message: "An Error Occured: " + err.Error(),
		}, nil
	}

	tokenMessage := &pb.String{Value: tokenString}
	anyToken, err := anypb.New(tokenMessage)
	if err != nil {
		return &pb.ApiResponse{
			Success: false,
			Message: "An Error Occured: " + err.Error(),
		}, nil
	}

	zap.L().Info("User Logged in. Token is: " + anyToken.String())

	return &pb.ApiResponse{
		Success: true,
		Message: "Login successful",
		Data:    anyToken,
	}, nil
}

func (h *Handlers) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.ApiResponse, error) {
	if req.Email == "" || req.Username == "" || req.Password == "" || req.Name == "" || req.SecurityQuestion == "" || req.SecurityAnswer == "" || req.Gender == "" || req.DateOfBirth == "" {
		return &pb.ApiResponse{Success: false, Message: "All fields must be filled."}, nil
	}

	if len(req.Name) < 5 || !regexp.MustCompile(`^[A-Za-z\s]+$`).MatchString(req.Name) {
		return &pb.ApiResponse{Success: false, Message: "Name must be more than 4 characters and contain only letters and spaces."}, nil
	}

	var existingUser models.User
	if err := h.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return &pb.ApiResponse{Success: false, Message: "Username is already taken."}, nil
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.com$`).MatchString(req.Email) {
		return &pb.ApiResponse{Success: false, Message: "Invalid email format. Must end with .com"}, nil
	}

	if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return &pb.ApiResponse{Success: false, Message: "Email is already registered."}, nil
	}

	if len(req.Password) < 8 ||
		!regexp.MustCompile(`[A-Z]`).MatchString(req.Password) ||
		!regexp.MustCompile(`[a-z]`).MatchString(req.Password) ||
		!regexp.MustCompile(`[0-9]`).MatchString(req.Password) ||
		!regexp.MustCompile(`[!@#~$%^&*()+|_]`).MatchString(req.Password) {
		return &pb.ApiResponse{Success: false, Message: "Password must be at least 8 characters long and include uppercase, lowercase, number, and special character."}, nil
	}

	if req.Gender != "male" && req.Gender != "female" {
		return &pb.ApiResponse{Success: false, Message: "Gender must be 'male' or 'female'."}, nil
	}

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return &pb.ApiResponse{Success: false, Message: "Invalid date of birth format. Use YYYY-MM-DD."}, nil
	}
	if time.Since(dob).Hours() < 13*365*24 {
		return &pb.ApiResponse{Success: false, Message: "You must be at least 13 years old to register."}, nil
	}

	validQuestions := map[string]bool{
		"What was the name of your first pet?":    true,
		"What city were you born in?":             true,
		"What is your favorite video game?":       true,
		"What was the name of your first school?": true,
		"What was your childhood nickname?":       true,
	}
	if !validQuestions[req.SecurityQuestion] {
		return &pb.ApiResponse{Success: false, Message: "Invalid security question selected."}, nil
	}

	// if err := verifyRecaptcha(req.RecaptchaToken); err != nil {
	// 	return &pb.ApiResponse{Success: false, Message: "reCAPTCHA verification failed."}, nil
	// }

	h.RequestVerificationCode(ctx, &pb.VerificationRequest{
		Email: req.Email,
	})

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
		UserId:           uuid.New().String(),
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
		WantsNewsletter:  req.WantsNewsletter,
		IsVerified:       false,
		IsDeactivated:    true,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		return &pb.ApiResponse{Success: false, Message: "Failed to register user."}, err
	}

	return &pb.ApiResponse{Success: true, Message: "User registered successfully. Please verify your email."}, nil
}
