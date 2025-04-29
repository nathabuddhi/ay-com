package supabase

import (
	"fmt"
	"io"
	"os"

	storage_go "github.com/supabase-community/storage-go"
	"go.uber.org/zap"
)

var Client *storage_go.Client

func InitSupabase() {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")

	storageURL := url + "/storage/v1"

	Client = storage_go.NewClient(storageURL, key, nil)
}

func UploadAvatarAndBanner(userId string, avatarFile io.Reader, bannerFile io.Reader) error {
	avatarBucket := "avatars"
	bannerBucket := "banners"

	avatarPath := fmt.Sprintf("%s.png", userId)
	if err := uploadFile(avatarBucket, avatarPath, avatarFile); err != nil {
		zap.L().Error("Failed to upload avatar", zap.Error(err))
		return err
	}

	bannerPath := fmt.Sprintf("%s.png", userId)
	if err := uploadFile(bannerBucket, bannerPath, bannerFile); err != nil {
		zap.L().Error("Failed to upload banner", zap.Error(err))
		return err
	}

	return nil
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
