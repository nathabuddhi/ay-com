package handlers

import (
	"fmt"
	"io"
	"strings"
	"sync"

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

func UploadThreadMedia(threadId string, threadFile io.Reader) error {
	imageData, err := readFileToBytes(threadFile)
	if err != nil {
		return fmt.Errorf("failed to read thread file: %w", err)
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

	bannerReq := &pb.UploadImageRequest{
		TypeId:     threadId,
		Image:      imageData,
		UploadType: "banner",
		ImageType:  fileExtension,
	}

	if _, err := client.Media_UploadMedia(ctx, bannerReq); err != nil {
		zap.L().Error("Failed to upload thread media", zap.Error(err))
		return fmt.Errorf("failed to upload thread media: %w", err)
	}

	return nil
}

func UploadMessageMedia(threadId string, threadFile io.Reader) error {
	imageData, err := readFileToBytes(threadFile)
	if err != nil {
		return fmt.Errorf("failed to read thread file: %w", err)
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

	bannerReq := &pb.UploadImageRequest{
		TypeId:     threadId,
		Image:      imageData,
		UploadType: "banner",
		ImageType:  fileExtension,
	}

	if _, err := client.Media_UploadMedia(ctx, bannerReq); err != nil {
		zap.L().Error("Failed to upload thread media", zap.Error(err))
		return fmt.Errorf("failed to upload thread media: %w", err)
	}

	return nil
}
