package server

import (
	"bytes"
	"context"

	pb "github.com/nathabuddhi/ay-com/backend/service-media/proto/media"
	"github.com/nathabuddhi/ay-com/backend/service-media/supabase"
	"go.uber.org/zap"
)

type MediaServer struct {
	pb.UnimplementedMediaServiceServer
}

func NewMediaServer() *MediaServer {
	return &MediaServer{}
}

func (s *MediaServer) Media_UploadMedia(ctx context.Context, req *pb.UploadImageRequest) (*pb.ApiResponseMedia, error) {
	if req.UploadType == "avatar" {
		path, err := supabase.UploadAvatar(req.TypeId, bytes.NewReader(req.Image))
		if err != nil {
			zap.L().Error("Failed to upload avatar", zap.Error(err))
			return &pb.ApiResponseMedia{
				Success: false,
				Message: "Failed to upload avatar.",
			}, nil
		}

		return &pb.ApiResponseMedia{
			Success: true,
			Message: "Avatar uploaded successfully.",
			Url:     path,
		}, nil
	}

	if req.UploadType == "banner" {
		path, err := supabase.UploadBanner(req.TypeId, bytes.NewReader(req.Image))
		if err != nil {
			zap.L().Error("Failed to upload banner", zap.Error(err))
			return &pb.ApiResponseMedia{
				Success: false,
				Message: "Failed to upload banner.",
			}, nil
		}

		return &pb.ApiResponseMedia{
			Success: true,
			Message: "Banner uploaded successfully.",
			Url:     path,
		}, nil
	}

	if req.UploadType == "thread" {
		path, err := supabase.UploadThreadMedia(req.TypeId, bytes.NewReader(req.Image))
		if err != nil {
			zap.L().Error("Failed to upload thread media", zap.Error(err))
			return &pb.ApiResponseMedia{
				Success: false,
				Message: "Failed to upload thread media.",
			}, nil
		}

		return &pb.ApiResponseMedia{
			Success: true,
			Message: "Thread Media uploaded successfully.",
			Url:     path,
		}, nil
	}

	if req.UploadType == "message" {
		// path, err := supabase.UploadMessageMedia(req.TypeId, bytes.NewReader(req.Image))
		// if err != nil {
		// 	zap.L().Error("Failed to upload message media", zap.Error(err))
		// 	return &pb.ApiResponseMedia{
		// 		Success: false,
		// 		Message: "Failed to upload message media.",
		// 	}, nil
		// }

		// return &pb.ApiResponseMedia{
		// 	Success: true,
		// 	Message: "Message Media uploaded successfully.",
		// 	Url:     path,
		// }, nil
	}

	return &pb.ApiResponseMedia{
		Success: false,
		Message: "Invalid upload type.",
	}, nil
}
