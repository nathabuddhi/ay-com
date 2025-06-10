package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/middleware"
	pb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/media"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	mediaServiceConn *grpc.ClientConn
	mediaServiceOnce sync.Once
)

func getMediaServiceConn() *grpc.ClientConn {
	mediaServiceOnce.Do(func() {
		var err error
		mediaServiceConn, err = grpc.NewClient(MEDIA_SERVICE_PATH, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			zap.L().Fatal("Failed to connect to media service: " + err.Error())
		}
	})
	return mediaServiceConn
}

func readFileToBytes(file io.Reader) ([]byte, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		zap.L().Error("Failed to read file", zap.Error(err))
		return nil, err
	}
	return data, nil
}

func UploadAvatar(userId string, avatarFile io.Reader) error {
	imageData, err := readFileToBytes(avatarFile)
	if err != nil {
		return fmt.Errorf("failed to read avatar file: %w", err)
	}

	conn := getMediaServiceConn()
	client := pb.NewMediaServiceClient(conn)
	ctx, cancel := createContext()
	defer cancel()

	avatarReq := &pb.UploadImageRequest{
		TypeId:     userId,
		Image:      imageData,
		UploadType: "avatar",
		ImageType:  "png",
	}

	if _, err := client.Media_UploadMedia(ctx, avatarReq); err != nil {
		zap.L().Error("Failed to upload avatar", zap.Error(err))
		return fmt.Errorf("failed to upload avatar: %w", err)
	}

	return nil
}

func UploadBanner(userId string, bannerFile io.Reader) error {
	imageData, err := readFileToBytes(bannerFile)
	if err != nil {
		return fmt.Errorf("failed to read banner file: %w", err)
	}

	conn := getMediaServiceConn()
	client := pb.NewMediaServiceClient(conn)
	ctx, cancel := createContext()
	defer cancel()

	bannerReq := &pb.UploadImageRequest{
		TypeId:     userId,
		Image:      imageData,
		UploadType: "banner",
		ImageType:  "png",
	}

	if _, err := client.Media_UploadMedia(ctx, bannerReq); err != nil {
		zap.L().Error("Failed to upload banner", zap.Error(err))
		return fmt.Errorf("failed to upload banner: %w", err)
	}

	return nil
}

func UploadThreadMedia(threadId string, threadFile io.Reader) (string, error) {
	imageData, err := readFileToBytes(threadFile)
	if err != nil {
		return "", fmt.Errorf("failed to read thread file: %w", err)
	}

	conn := getMediaServiceConn()
	client := pb.NewMediaServiceClient(conn)
	ctx, cancel := createContext()
	defer cancel()

	fileExtension := "png"

	if file, ok := threadFile.(interface{ Name() string }); ok {
		fileName := file.Name()
		if dotIndex := strings.LastIndex(fileName, "."); dotIndex != -1 && dotIndex < len(fileName)-1 {
			fileExtension = fileName[dotIndex+1:]
		}
	}
	if file, ok := threadFile.(interface{ Name() string }); ok {
		fileName := file.Name()
		if dotIndex := strings.LastIndex(fileName, "."); dotIndex != -1 && dotIndex < len(fileName)-1 {
			fileExtension = fileName[dotIndex+1:]
		}
	}

	bannerReq := &pb.UploadImageRequest{
		TypeId:     threadId,
		Image:      imageData,
		UploadType: "thread",
		ImageType:  fileExtension,
	}

	if resp, err := client.Media_UploadMedia(ctx, bannerReq); err != nil {
		zap.L().Error("Failed to upload thread media", zap.Error(err))
		return "", fmt.Errorf("failed to upload thread media: %w", err)
	} else {
		return resp.Url, nil
	}

}

func UploadMessageMedia(messageId string, messageFile io.Reader) error {
	imageData, err := readFileToBytes(messageFile)
	if err != nil {
		return fmt.Errorf("failed to read thread file: %w", err)
	}

	conn := getMediaServiceConn()
	client := pb.NewMediaServiceClient(conn)
	ctx, cancel := createContext()
	defer cancel()

	fileExtension := "png"
	if file, ok := messageFile.(interface{ Name() string }); ok {
		fileName := file.Name()
		if dotIndex := strings.LastIndex(fileName, "."); dotIndex != -1 && dotIndex < len(fileName)-1 {
			fileExtension = fileName[dotIndex+1:]
		}
	}

	bannerReq := &pb.UploadImageRequest{
		TypeId:     messageId,
		Image:      imageData,
		UploadType: "message",
		ImageType:  fileExtension,
	}

	if _, err := client.Media_UploadMedia(ctx, bannerReq); err != nil {
		zap.L().Error("Failed to upload thread media", zap.Error(err))
		return fmt.Errorf("failed to upload thread media: %w", err)
	}

	return nil
}

func UploadVerifyRequestImage(faceFile io.Reader, fileExtension string) (string, error) {
	imageData, err := readFileToBytes(faceFile)
	if err != nil {
		return "", fmt.Errorf("failed to read face file: %w", err)
	}

	conn := getMediaServiceConn()
	client := pb.NewMediaServiceClient(conn)
	ctx, cancel := createContext()
	defer cancel()

	generatedId := uuid.New().String()

	bannerReq := &pb.UploadImageRequest{
		TypeId:     generatedId,
		Image:      imageData,
		UploadType: "verification",
		ImageType:  fileExtension,
	}

	if _, err := client.Media_UploadMedia(ctx, bannerReq); err != nil {
		zap.L().Error("Failed to upload thread media", zap.Error(err))
		return "", fmt.Errorf("failed to upload thread media: %w", err)
	}

	return generatedId + fileExtension, nil
}

func User_ChangeAvatar(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Change Avatar is called.")

	conn := getUserServiceConn()
	client := pb.NewMediaServiceClient(conn)

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

	if len(avatarHeader.Filename) < 4 || strings.ToLower(avatarHeader.Filename[len(avatarHeader.Filename)-4:]) != ".png" {
		returnErrorResponse(w, "Avatar file must be a .png file")
		return
	}

	imageData, err := readFileToBytes(avatarFile)
	if err != nil {
		returnErrorResponse(w, "Failed to read avatar file.")
		return
	}

	uploadImageReq := &pb.UploadImageRequest{
		TypeId:     r.Context().Value(middleware.UserIdKey).(string),
		Image:      imageData,
		UploadType: "avatar",
		ImageType:  "png",
	}
	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Media_UploadMedia(ctx, uploadImageReq)
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func User_ChangeBanner(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("User Change Banner is called.")

	conn := getUserServiceConn()
	client := pb.NewMediaServiceClient(conn)

	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		zap.L().Error("Failed to parse multipart form", zap.Error(err))
		returnErrorResponse(w, "Failed to parse multipart form")
		return
	}

	bannerFile, bannerHeader, err := r.FormFile("banner")
	if err != nil {
		zap.L().Error("Failed to get banner file", zap.Error(err))
		returnErrorResponse(w, "Failed to get banner file")
		return
	}
	defer bannerFile.Close()

	if len(bannerHeader.Filename) < 4 || strings.ToLower(bannerHeader.Filename[len(bannerHeader.Filename)-4:]) != ".png" {
		returnErrorResponse(w, "Banner file must be a .png file")
		return
	}

	imageData, err := readFileToBytes(bannerFile)
	if err != nil {
		returnErrorResponse(w, "Failed to read banner file.")
		return
	}

	uploadImageReq := &pb.UploadImageRequest{
		TypeId:     r.Context().Value(middleware.UserIdKey).(string),
		Image:      imageData,
		UploadType: "banner",
		ImageType:  "png",
	}
	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Media_UploadMedia(ctx, uploadImageReq)
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
