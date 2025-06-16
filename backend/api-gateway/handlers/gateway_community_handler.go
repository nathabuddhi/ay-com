package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	categoryCountStr := r.FormValue("category_count")
	var categoryCount int32
	if categoryCountStr != "" {
		var parsedCount int
		_, err := fmt.Sscanf(categoryCountStr, "%d", &parsedCount)
		if err != nil {
			zap.L().Error("Invalid category_count value", zap.Error(err))
			returnErrorResponse(w, "Invalid category_count value")
			return
		}
		categoryCount = int32(parsedCount)
	}

	if categoryCount < 1 {
		zap.L().Error("Category count must be at least 1")
		returnErrorResponse(w, "Category count must be at least 1")
		return
	}

	req.Categories = make([]string, categoryCount)
	for i := 0; i < int(categoryCount); i++ {
		category := r.FormValue(fmt.Sprintf("category_%d", i+1))
		if category == "" {
			zap.L().Error("Category cannot be empty", zap.Int("category_index", i+1))
			returnErrorResponse(w, fmt.Sprintf("Category %d cannot be empty", i+1))
			return
		}
		req.Categories[i] = category
	}

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

func Community_GetUserCommunities(w http.ResponseWriter, r *http.Request) {
	conn := getCommunityServiceConn()
	client := pb.NewCommunityServiceClient(conn)

	req := &pb.StringCommunity{}
	req.Value = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Community_GetUserCommunities(ctx, req)

	processCommunityResponseWithPayload[pb.GetCommunitiesResponse](resp, err, w)
}

func Community_GetCategories(w http.ResponseWriter, r *http.Request) {
	conn := getCommunityServiceConn()
	client := pb.NewCommunityServiceClient(conn)

	req := &pb.StringCommunity{}
	req.Value = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Community_GetCategories(ctx, req)

	if err != nil {
		zap.L().Error("Error forwarding request", zap.Error(err))
		returnErrorResponse(w, "Error forwarding request: "+err.Error())
		return
	}

	if resp == nil || resp.Categories == nil {
		zap.L().Error("Received nil response from community service")
		returnErrorResponse(w, "Received nil response from community service")
		return
	}

	response := &types.ApiResponse{
		Success: true,
		Message: "Successfully retrieved categories.",
		Payload: &pb.GetCategoriesResponse{
			Categories: resp.Categories,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Community_GetCommunityById(w http.ResponseWriter, r *http.Request) {
	conn := getCommunityServiceConn()
	client := pb.NewCommunityServiceClient(conn)

	vars := mux.Vars(r)
	community_id := vars["id"]
	if community_id == "" {
		zap.L().Error("Community id parameter missing")
		returnErrorResponse(w, "Community id is required")
		return
	}

	if !checkRedisData("getcommunity/"+community_id, w) {
		req := &pb.GeneralCommunityRequest{}

		req.CommunityId = community_id
		req.UserId = r.Context().Value(middleware.UserIdKey).(string)

		ctx, cancel := createContext()
		defer cancel()

		resp, err := client.Community_GetCommunityById(ctx, req)

		processCommunityResponseWithPayload[pb.Community](resp, err, w)
	}
}

func Community_GetAllCommunities(w http.ResponseWriter, r *http.Request) {
	conn := getCommunityServiceConn()
	client := pb.NewCommunityServiceClient(conn)

	req := &pb.StringCommunity{}
	req.Value = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Community_GetAllCommunities(ctx, req)

	processCommunityResponseWithPayload[pb.GetCommunitiesResponse](resp, err, w)
}

func Community_GetUserPendingCommunities(w http.ResponseWriter, r *http.Request) {
	conn := getCommunityServiceConn()
	client := pb.NewCommunityServiceClient(conn)

	req := &pb.StringCommunity{}
	req.Value = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Community_GetUserPendingCommunities(ctx, req)

	processCommunityResponseWithPayload[pb.GetCommunitiesResponse](resp, err, w)
}

func Community_GetCommunityMembers(w http.ResponseWriter, r *http.Request) {
	req, client := processCommunityRequest[pb.StringCommunity](r, w)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Community_GetCommunityMembers(ctx, req)

	processCommunityResponseWithPayload[pb.GetMembersResponse](resp, err, w)
}

func Community_ApproveMember(w http.ResponseWriter, r *http.Request) {
	req, client := processCommunityRequest[pb.GeneralModeratorRequest](r, w)
	req.ModeratorId = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Community_ApproveJoinRequest(ctx, req)

	processCommunityResponseWithoutPayload(resp, err, w)
}

func Community_JoinCommunity(w http.ResponseWriter, r *http.Request) {
	req, client := processCommunityRequest[pb.GeneralCommunityRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Community_JoinCommunity(ctx, req)

	processCommunityResponseWithoutPayload(resp, err, w)
}

func Community_DenyMember(w http.ResponseWriter, r *http.Request) {
	req, client := processCommunityRequest[pb.GeneralModeratorRequest](r, w)
	req.ModeratorId = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Community_DenyJoinRequest(ctx, req)

	processCommunityResponseWithoutPayload(resp, err, w)
}

func Community_GetJoinRequests(w http.ResponseWriter, r *http.Request) {
	req, client := processCommunityRequest[pb.StringCommunity](r, w)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Community_GetJoinRequests(ctx, req)

	processCommunityResponseWithoutPayload(resp, err, w)
}
