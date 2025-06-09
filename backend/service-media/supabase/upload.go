package supabase

import (
	"bytes"
	"fmt"
	"io"

	pb "github.com/nathabuddhi/ay-com/backend/service-media/proto/media"
	"go.uber.org/zap"
)

func UploadAvatar(userId string, avatarFile io.Reader) (string, error) {
	zap.L().Info("Uploading Avatar", zap.String("user_id", userId))

	avatarBucket := "avatars"

	avatarPath := fmt.Sprintf("%s.png", userId)
	if err := uploadFile(avatarBucket, avatarPath, avatarFile); err != nil {
		zap.L().Error("Failed to upload avatar", zap.Error(err))
		return "", err
	}

	return avatarPath, nil
}

func UploadVerificationMedia(req *pb.UploadImageRequest) (string, error) {
	zap.L().Info("Uploading verification request", zap.String("id", req.TypeId))

	bucket := "verificationrequest"

	file := bytes.NewReader(req.Image)

	path := fmt.Sprintf("%s.%s", req.TypeId, req.ImageType)
	if err := uploadFile(bucket, path, file); err != nil {
		zap.L().Error("Failed to upload avatar", zap.Error(err))
		return "", err
	}

	return path, nil
}

func UploadBanner(userId string, bannerFile io.Reader) (string, error) {
	zap.L().Info("Uploading Banner", zap.String("user_id", userId))

	bannerBucket := "banners"

	bannerPath := fmt.Sprintf("%s.png", userId)
	if err := uploadFile(bannerBucket, bannerPath, bannerFile); err != nil {
		zap.L().Error("Failed to upload banner", zap.Error(err))
		return "", err
	}

	return bannerPath, nil
}

func UploadThreadMedia(thread_id string, file io.Reader) (string, error) {
	zap.L().Info("Uploading Thread Media", zap.String("thread_id", thread_id))

	mediaBucket := "threads"

	var mediaPath string
	suffix := 0

	for {
		mediaPath = fmt.Sprintf("%s/%d.png", thread_id, suffix)
		if !fileExists(mediaBucket, mediaPath) {
			break
		}
		suffix++
	}

	err := uploadFile(mediaBucket, mediaPath, file)
	if err != nil {
		zap.L().Error("Failed to upload thread media", zap.Error(err))
		return "", err
	}

	return mediaPath, nil
}

func fileExists(bucket, filePath string) bool {
	_, err := Client.DownloadFile(bucket, filePath)
	if err != nil {
		zap.L().Debug("File does not exist or download failed", zap.String("filePath", filePath), zap.Error(err))
		return false
	}
	return true
}

func deleteFile(bucket string, mediaPath string) bool {
	if _, err := Client.RemoveFile(bucket, []string{mediaPath}); err != nil {
		return false
	}
	return true
}

func uploadFile(bucketName, filePath string, file io.Reader) error {
	if fileExists(bucketName, filePath) {
		if !deleteFile(bucketName, filePath) {
			return fmt.Errorf("failed to delete existing file: %s", filePath)
		}
	}

	_, err := Client.UploadFile(bucketName, filePath, file)
	if err != nil {
		zap.L().Error("Failed to upload file to Supabase: " + err.Error())
		return err
	}

	zap.L().Info("Successfully uploaded file to Supabase ", zap.String("path", filePath))
	return nil
}
