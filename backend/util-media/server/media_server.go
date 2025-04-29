package server

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	pb "github.com/nathabuddhi/ay-com/backend/util-media/proto"
	"github.com/nathabuddhi/ay-com/backend/util-media/supabase"
)

type MediaServer struct {
	pb.UnimplementedMediaServiceServer
}

func NewMediaServer() *MediaServer {
	return &MediaServer{}
}

func (s *MediaServer) UploadMedia(ctx context.Context, req *pb.UploadMediaRequest) (*pb.UploadMediaResponse, error) {
	var bucket string
	var path string

	switch req.UsageType {
	case "avatar":
		bucket = "avatars"
		extension := getFileExtension(req.FileName)
		path = fmt.Sprintf("%s.%s", req.OwnerId, extension)
	case "banner":
		bucket = "banners"
		extension := getFileExtension(req.FileName)
		path = fmt.Sprintf("%s.%s", req.OwnerId, extension)
	case "thread":
		bucket = "threads"
		path = fmt.Sprintf("%s/%s", req.OwnerId, req.FileName)
	case "message":
		bucket = "messages"
		path = fmt.Sprintf("%s/%s", req.OwnerId, req.FileName)
	default:
		return &pb.UploadMediaResponse{Success: false}, fmt.Errorf("invalid usage type")
	}

	fileReader := bytes.NewReader(req.File)

	_, err := supabase.Client.UploadFile(bucket, path, fileReader)
	if err != nil {
		return &pb.UploadMediaResponse{Success: false}, err
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", os.Getenv("SUPABASE_URL"), bucket, path)

	return &pb.UploadMediaResponse{
		Success: true,
		Url:     publicURL,
	}, nil
}

func (s *MediaServer) DeleteMedia(ctx context.Context, req *pb.DeleteMediaRequest) (*pb.DeleteMediaResponse, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	pathPrefix := fmt.Sprintf("%s/storage/v1/object/public/", supabaseURL)

	filePath := req.FileUrl[len(pathPrefix):]

	slashIndex := findNth(filePath, '/', 1)
	if slashIndex == -1 {
		return &pb.DeleteMediaResponse{Success: false}, fmt.Errorf("invalid file url")
	}

	bucket := filePath[:slashIndex]
	path := filePath[slashIndex+1:]

	_, err := supabase.Client.RemoveFile(bucket, []string{path})
	if err != nil {
		return &pb.DeleteMediaResponse{Success: false}, err
	}

	return &pb.DeleteMediaResponse{
		Success: true,
	}, nil
}

func getFileExtension(filename string) string {
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")
	return ext
}

func findNth(s string, sep byte, n int) int {
	count := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			count++
			if count == n {
				return i
			}
		}
	}
	return -1
}
