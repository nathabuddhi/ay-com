package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/middleware"
	pb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/community"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/types"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

var (
	communityServiceConn *grpc.ClientConn
	communityServiceOnce sync.Once
)

func getCommunityServiceConn() *grpc.ClientConn {
	communityServiceOnce.Do(func() {
		var err error
		communityServiceConn, err = grpc.NewClient(COMMUNITY_SERVICE_PATH, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			zap.L().Fatal("Failed to connect to community service: " + err.Error())
		}
	})
	return communityServiceConn
}

func processCommunityResponseWithoutPayload(resp *pb.ApiResponseCommunity, err error, w http.ResponseWriter) {
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	response := &types.ApiResponse{
		Success: resp.Success,
		Message: resp.Message,
		Payload: nil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func processCommunityResponseWithPayload[T any](resp *pb.ApiResponseCommunity, err error, w http.ResponseWriter) {
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	decodedObject := new(T)
	if _, ok := any(decodedObject).(proto.Message); !ok {
		zap.L().Error("Type does not implement proto.Message", zap.String("type", fmt.Sprintf("%T", decodedObject)))
		returnErrorResponse(w, "Failed to decode response data.")
		return
	}

	if err := proto.Unmarshal(resp.Data.GetValue(), any(decodedObject).(proto.Message)); err != nil {
		zap.L().Error("Failed to unmarshal protobuf", zap.Error(err))
		returnErrorResponse(w, "Failed to decode response data.")
		return
	}

	response := &types.ApiResponse{
		Success: resp.Success,
		Message: resp.Message,
		Payload: decodedObject,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func processCommunityRequest[T any](r *http.Request, w http.ResponseWriter) (resp *T, client pb.CommunityServiceClient) {
	conn := getCommunityServiceConn()
	client = pb.NewCommunityServiceClient(conn)

	var req T
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to encode API request", zap.Error(err))
		returnErrorResponse(w, "Failed to encode API request: "+err.Error())
	}
	return &req, client
}

func Community_CreateCommunity(w http.ResponseWriter, r *http.Request) {
	conn := getCommunityServiceConn()
	client := pb.NewCommunityServiceClient(conn)

	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		zap.L().Error("Failed to parse multipart form", zap.Error(err))
		returnErrorResponse(w, "Failed to parse multipart form")
		return
	}

	iconFile, iconHeader, err := r.FormFile("icon")
	if err != nil {
		zap.L().Error("Failed to get icon file", zap.Error(err))
		returnErrorResponse(w, "Failed to get icon file")
		return
	}
	defer iconFile.Close()

	bannerFile, bannerHeader, err := r.FormFile("banner")
	if err != nil {
		zap.L().Error("Failed to get banner file", zap.Error(err))
		returnErrorResponse(w, "Failed to get banner file")
		return
	}
	defer bannerFile.Close()

	generatedId := uuid.New().String()

	iconPath, err := UploadCommunityIcon(iconFile, iconHeader, generatedId)
	if err != nil {
		zap.L().Error("Error uploading icon.", zap.Error(err))
		returnErrorResponse(w, "Failed to upload icon file.")
		return
	}

	bannerPath, err := UploadCommunityBanner(bannerFile, bannerHeader, generatedId)
	if err != nil {
		zap.L().Error("Error uploading banner.", zap.Error(err))
		returnErrorResponse(w, "Failed to upload banner file.")
		return
	}

	req := &pb.CreateCommunityRequest{}
	req.Category = r.FormValue("category")
	req.CommunityName = r.FormValue("name")
	req.Description = r.FormValue("description")
	req.Rules = r.FormValue("rules")
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	req.BannerImage = bannerPath
	req.IconImage = iconPath
	req.CommunityId = generatedId

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Community_CreateCommunity(ctx, req)
	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	processCommunityResponseWithoutPayload(resp, err, w)
}
