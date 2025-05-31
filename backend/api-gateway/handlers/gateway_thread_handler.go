package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/middleware"
	pb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/thread"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/types"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

var (
	threadServiceConn *grpc.ClientConn
	threadServiceOnce sync.Once
)

func getThreadServiceConn() *grpc.ClientConn {
	threadServiceOnce.Do(func() {
		var err error
		threadServiceConn, err = grpc.NewClient(THREAD_SERVICE_PATH, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			zap.L().Fatal("Failed to connect to thread service: " + err.Error())
		}
	})
	return threadServiceConn
}

func processThreadResponseWithoutPayload(resp *pb.ApiResponseThread, err error, w http.ResponseWriter) {
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

func processThreadResponseWithPayload[T any](resp *pb.ApiResponseThread, err error, w http.ResponseWriter) {
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

func processThreadRequest[T any](r *http.Request, w http.ResponseWriter) (resp *T, client pb.ThreadServiceClient) {
	conn := getThreadServiceConn()
	client = pb.NewThreadServiceClient(conn)

	var req T
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("Failed to encode API request", zap.Error(err))
		returnErrorResponse(w, "Failed to encode API request: "+err.Error())
	}
	return &req, client
}

func Thread_CreateThread(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (CreateThread) is called.")

	conn := getThreadServiceConn()
	client := pb.NewThreadServiceClient(conn)

	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		zap.L().Error("Failed to parse multipart form", zap.Error(err))
		returnErrorResponse(w, "Failed to parse multipart form")
		return
	}

	req := &pb.PostThread{}
	mediaCount, err := strconv.Atoi(r.FormValue("media_count"))
	if err != nil {
		returnErrorResponse(w, "Invalid media.")
		return
	}
	req.MediaCount = int32(mediaCount)

	req.MediaUrls = make([]string, req.MediaCount)
	folderId := uuid.New().String()

	for i := int32(0); i < req.MediaCount; i++ {
		mediaFile, mediaHeader, err := r.FormFile("media_" + strconv.Itoa(int(i)))
		if err != nil {
			zap.L().Error("Failed to get media file", zap.Error(err))
			returnErrorResponse(w, "Failed to get media file")
			return
		}
		defer mediaFile.Close()

		ext := strings.ToLower(mediaHeader.Filename[strings.LastIndex(mediaHeader.Filename, "."):])
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".gif" && ext != ".mp4" {
			returnErrorResponse(w, "Media file must be a .png, .jpg, .jpeg, .gif, or .mp4 file")
			return
		}

		url, err := UploadThreadMedia(folderId, mediaFile)
		if err != nil {
			returnErrorResponse(w, "Failed to upload media files.")
			return
		}
		req.MediaUrls[i] = url
	}

	pollCount, err := strconv.Atoi(r.FormValue("poll_count"))
	if err != nil {
		returnErrorResponse(w, "Invalid media.")
		return
	}
	req.PollCount = int32(pollCount)

	req.PollOptions = make([]string, req.PollCount)

	for i := int32(0); i < req.PollCount; i++ {
		pollOption := r.FormValue("poll_" + strconv.Itoa(int(i)))

		req.PollOptions[i] = pollOption
	}

	req.UserId = r.Context().Value(middleware.UserIdKey).(string)
	req.Title = r.FormValue("title")
	req.Content = r.FormValue("content")
	req.Category = r.FormValue("category")
	req.CommunityId = r.FormValue("community_id")
	req.ReplyPermission = r.FormValue("reply_permission")

	if r.FormValue("is_private") == "true" {
		req.IsPrivate = true
	} else {
		req.IsPrivate = false
	}

	if r.FormValue("is_advertisement") == "true" {
		req.IsAdvertisement = true
	} else {
		req.IsAdvertisement = false
	}

	if r.FormValue("is_scheduled") == "true" {
		req.IsScheduled = true
		req.ScheduledAt = r.FormValue("scheduled_at")
	} else {
		req.IsAdvertisement = false
	}

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_CreateThread(ctx, req)
	processThreadResponseWithoutPayload(resp, err, w)
}

func Thread_GetAllThreads(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (GetAllThreads) is called.")

	req, client := processThreadRequest[pb.GetAllThreadsRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_GetAllThreads(ctx, req)
	processThreadResponseWithPayload[pb.GetThreadsResponse](resp, err, w)
}

func Thread_SearchThreads(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (SearchThreads) is called.")

	req, client := processThreadRequest[pb.GetAllThreadsRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_SearchThreads(ctx, req)
	processThreadResponseWithPayload[pb.GetThreadsResponse](resp, err, w)
}

func Thread_DeleteThread(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (DeleteThread) is called.")

	req, client := processThreadRequest[pb.DeleteThreadRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_DeleteThread(ctx, req)
	processThreadResponseWithoutPayload(resp, err, w)
}

func Thread_PinThread(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (PinThread) is called.")

	req, client := processThreadRequest[pb.GeneralThreadRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_PinThread(ctx, req)
	processThreadResponseWithoutPayload(resp, err, w)
}

func Thread_VoteThread(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (VoteThread) is called.")

	req, client := processThreadRequest[pb.SubmitVote](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_VoteThread(ctx, req)
	processThreadResponseWithoutPayload(resp, err, w)
}

func Thread_ToggleLike(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (Toggle Like) is called.")

	req, client := processThreadRequest[pb.GeneralThreadRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_ToggleLike(ctx, req)
	processThreadResponseWithoutPayload(resp, err, w)
}

func Thread_ToggleBookmark(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (Toggle BookMark) is called.")

	req, client := processThreadRequest[pb.GeneralThreadRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_ToggleBookmark(ctx, req)
	processThreadResponseWithoutPayload(resp, err, w)
}

func Thread_ToggleRepost(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (ToggleRepost) is called.")

	req, client := processThreadRequest[pb.RepostRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_ToggleRepost(ctx, req)
	processThreadResponseWithoutPayload(resp, err, w)
}

func Thread_ReplyThread(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (ReplyThread) is called.")

	req, client := processThreadRequest[pb.GeneralThreadRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_ReplyThread(ctx, req)
	processThreadResponseWithoutPayload(resp, err, w)
}

func Thread_DeleteReply(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (DeleteReply) is called.")

	req, client := processThreadRequest[pb.GeneralThreadRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_DeleteReply(ctx, req)
	processThreadResponseWithoutPayload(resp, err, w)
}
