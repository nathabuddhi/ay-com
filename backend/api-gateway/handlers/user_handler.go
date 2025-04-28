package handlers

import (
	"context"
	"net/http"
	"sync"

	pb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/user"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	forwardRequest[pb.LoginRequest, pb.String](w, r, func(ctx context.Context, in *pb.LoginRequest) (*pb.ApiResponse, error) {
		return client.User_Login(ctx, in)
	})
}

func User_Register(w http.ResponseWriter, r *http.Request) {
	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)
	forwardRequest[pb.RegisterRequest, pb.UserProfile](w, r, func(ctx context.Context, in *pb.RegisterRequest) (*pb.ApiResponse, error) {
		return client.User_Register(ctx, in)
	})
}

func User_GetProfile(w http.ResponseWriter, r *http.Request) {
	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)

	forwardRequest[pb.GetUserRequest, pb.UserProfile](w, r, func(ctx context.Context, in *pb.GetUserRequest) (*pb.ApiResponse, error) {
		return client.User_GetProfile(ctx, in)
	})
}

func User_RequestVerificationCode(w http.ResponseWriter, r *http.Request) {
	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)
	forwardRequest[pb.VerificationRequest, pb.String](w, r, func(ctx context.Context, in *pb.VerificationRequest) (*pb.ApiResponse, error) {
		return client.RequestVerificationCode(ctx, in)
	})
}

func User_ValidateVerificationCode(w http.ResponseWriter, r *http.Request) {
	conn := getUserServiceConn()
	client := pb.NewUserServiceClient(conn)
	forwardRequest[pb.ValidateCodeRequest, pb.String](w, r, func(ctx context.Context, in *pb.ValidateCodeRequest) (*pb.ApiResponse, error) {
		return client.ValidateVerificationCode(ctx, in)
	})
}
