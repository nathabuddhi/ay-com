package supabase

import (
	"fmt"
	"io"

	storage_go "github.com/supabase-community/storage-go"
	"go.uber.org/zap"
)

func UploadAvatar(userId string, avatarFile io.Reader) (string, error) {
	avatarBucket := "avatars"

	avatarPath := fmt.Sprintf("%s.png", userId)
	if err := uploadFile(avatarBucket, avatarPath, avatarFile); err != nil {
		zap.L().Error("Failed to upload avatar", zap.Error(err))
		return "", err
	}

	return avatarPath, nil
}

func UploadBanner(userId string, bannerFile io.Reader) (string, error) {
	bannerBucket := "banners"

	bannerPath := fmt.Sprintf("%s.png", userId)
	if err := uploadFile(bannerBucket, bannerPath, bannerFile); err != nil {
		zap.L().Error("Failed to upload banner", zap.Error(err))
		return "", err
	}

	return bannerPath, nil
}

func UploadThreadMedia(thread_id string, file io.Reader) (string, error) {
	mediaBucket := "threads"

	var mediaPath string
	suffix := 1

	for {
		mediaPath = fmt.Sprintf("%s/%d.png", thread_id, suffix)
		if !fileExists(mediaBucket, mediaPath) {
			break
		}
		suffix++
	}

	if err := uploadFile(mediaBucket, mediaPath, file); err != nil {
		zap.L().Error("Failed to upload thread media", zap.Error(err))
		return "", err
	}

	return mediaPath, nil
}

func fileExists(bucket, mediaPath string) bool {
	file, err := Client.ListFiles(bucket, mediaPath, storage_go.FileSearchOptions{})
	if err != nil {
		zap.L().Error("Error checking file existence", zap.Error(err))
		return true
	}
	return file != nil
}

func uploadFile(bucketName, filePath string, file io.Reader) error {
	_, err := Client.UploadFile(bucketName, filePath, file)
	if err != nil {
		zap.L().Error("Failed to upload file to Supabase: " + err.Error())
		return err
	}

	zap.L().Info("Successfully uploaded file to Supabase ", zap.String("path", filePath))
	return nil
}
