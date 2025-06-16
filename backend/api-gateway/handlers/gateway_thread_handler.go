package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/middleware"
	pb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/thread"
	userpb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/user"
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

func Thread_GetThreadById(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (GetThreadById) is called.")
	vars := mux.Vars(r)
	threadId := vars["id"]
	if !checkRedisData("getthread/"+threadId, w) {
		conn := getThreadServiceConn()
		client := pb.NewThreadServiceClient(conn)

		req := &pb.GeneralThreadRequest{}
		req.ThreadId = threadId
		req.UserId = r.Context().Value(middleware.UserIdKey).(string)
		ctx, cancel := createContext()
		defer cancel()
		resp, err := client.Thread_GetThreadById(ctx, req)
		processThreadResponseWithPayload[pb.GetThreadDetailResponse](resp, err, w)
	}
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
	if req.CommunityId != "" {
		req.IsCommunity = true
	}
	req.ReplyPermission = r.FormValue("reply_permission")
	req.ReplyTo = r.FormValue("reply_to")

	if req.ReplyTo != "" {
		ctx, cancel := createContext()
		defer cancel()

		getThreadReq := &pb.GeneralThreadRequest{}
		getThreadReq.ThreadId = req.ReplyTo
		resp, err := client.Thread_GetThreadById(ctx, getThreadReq)
		if err != nil || !resp.Success {
			returnErrorResponse(w, "Failed checking reply permissions: "+resp.Message)
			return
		} else {
			// AMBIL THREAD REPLY PERMISSION
			decodedObject := new(pb.GetThreadDetailResponse)
			if _, ok := any(decodedObject).(proto.Message); !ok {
				zap.L().Error("Type does not implement proto.Message", zap.String("type", fmt.Sprintf("%T", decodedObject)))
				returnErrorResponse(w, "Failed to decode response data while checking reply permissions.")
				return
			}

			if err := proto.Unmarshal(resp.Data.GetValue(), any(decodedObject).(proto.Message)); err != nil {
				zap.L().Error("Failed to unmarshal protobuf", zap.Error(err))
				returnErrorResponse(w, "Failed to decode response data while checking reply permissions.")
				return
			}

			if decodedObject.Thread.UserId == r.Context().Value(middleware.UserIdKey).(string) {
				// do nothing
			} else if decodedObject.Thread.ReplyPermission == "Accounts You Follow" {
				userConn := getUserServiceConn()
				userClient := userpb.NewUserServiceClient(userConn)

				ctx, cancel := createContext()
				defer cancel()

				userResp, err := userClient.User_IsUserFollowing(ctx, &userpb.IsUserFollowingRequest{
					FollowingId: r.Context().Value(middleware.UserIdKey).(string),
					FollowerId:  decodedObject.Thread.UserId,
				})

				if err != nil || !userResp.Value {
					returnErrorResponse(w, "You are not allowed to reply to this thread. You must be followed by the author of this thread to reply.")
					return
				}

			} else if decodedObject.Thread.ReplyPermission == "Verified Accounts" {
				userConn := getUserServiceConn()
				userClient := userpb.NewUserServiceClient(userConn)

				ctx, cancel := createContext()
				defer cancel()

				userResp, _ := userClient.User_GetUserId(ctx, &userpb.GetProfileRequest{
					UserId: r.Context().Value(middleware.UserIdKey).(string),
				})
				decodedObject := new(userpb.UserProfile)
				if _, ok := any(decodedObject).(proto.Message); !ok {
					zap.L().Error("Type does not implement proto.Message", zap.String("type", fmt.Sprintf("%T", decodedObject)))
					returnErrorResponse(w, "Failed to decode response data while checking reply permissions.")
					return
				}

				if err := proto.Unmarshal(userResp.Data.GetValue(), any(decodedObject).(proto.Message)); err != nil {
					zap.L().Error("Failed to unmarshal protobuf", zap.Error(err))
					returnErrorResponse(w, "Failed to decode response data while checking reply permissions.")
					return
				}
				zap.L().Info("Decoded user profile", zap.String("userId", decodedObject.Username), zap.Bool("isVerified", decodedObject.IsVerified))
				if !decodedObject.IsVerified {
					returnErrorResponse(w, "You are not allowed to reply to this thread. You must be a verified account to reply.")
					return
				}
			}
		}
	}

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
		req.IsScheduled = false
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

func Thread_GetFollowingThreads(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (GetFollowingThreads) is called.")

	var req pb.GetFollowingThreadRequest
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	req2, client2 := processUserRequest[userpb.GetAllFollowingRequest](r, w)
	req2.RequesterId = req.UserId
	req2.UserId = req.UserId
	ctx2, cancel2 := createContext()
	defer cancel2()
	resp2, err2 := client2.User_GetAllFollowing(ctx2, req2)

	if err2 != nil {
		zap.L().Error("Failed to get following list", zap.Error(err2))
		returnErrorResponse(w, "Failed to get following list: "+err2.Error())
		return
	}

	decodedObject := new(userpb.AllFollowingResponse)
	if _, ok := any(decodedObject).(proto.Message); !ok {
		zap.L().Error("Type does not implement proto.Message", zap.String("type", fmt.Sprintf("%T", decodedObject)))
		returnErrorResponse(w, "Failed to decode response data.")
		return
	}

	if err2 := proto.Unmarshal(resp2.Data.GetValue(), any(decodedObject).(proto.Message)); err2 != nil {
		zap.L().Error("Failed to unmarshal protobuf", zap.Error(err2))
		returnErrorResponse(w, "Failed to decode response data.")
		return
	}

	for _, v := range decodedObject.Following {
		req.FollowedIds = append(req.FollowedIds, v.Value)
	}

	ctx, cancel := createContext()
	defer cancel()

	conn := getThreadServiceConn()
	client := pb.NewThreadServiceClient(conn)

	resp, err := client.Thread_GetFollowingThreads(ctx, &req)

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

func Thread_GetLatestThreads(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (SearchThreads) is called.")

	req, client := processThreadRequest[pb.GetAllThreadsRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_GetLatestThreads(ctx, req)
	processThreadResponseWithPayload[pb.GetThreadsResponse](resp, err, w)
}

func Thread_DeleteThread(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (DeleteThread) is called.")

	req, client := processThreadRequest[pb.GeneralThreadRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_DeleteThread(ctx, req)
	processThreadResponseWithoutPayload(resp, err, w)
}

func Thread_TogglePinThread(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (PinThread) is called.")

	req, client := processThreadRequest[pb.GeneralThreadRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_TogglePinThread(ctx, req)
	processThreadResponseWithoutPayload(resp, err, w)
}

func Thread_VoteThread(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (VoteThread) is called.")

	req, client := processThreadRequest[pb.SubmitVote](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	getThreadReq := &pb.GeneralThreadRequest{}
	getThreadReq.ThreadId = req.ThreadId
	resp, err := client.Thread_GetThreadById(ctx, getThreadReq)
	if err != nil || !resp.Success {
		returnErrorResponse(w, "Failed checking voting permissions: "+resp.Message)
		return
	} else {
		// AMBIL THREAD REPLY PERMISSION
		decodedObject := new(pb.GetThreadDetailResponse)
		if _, ok := any(decodedObject).(proto.Message); !ok {
			zap.L().Error("Type does not implement proto.Message", zap.String("type", fmt.Sprintf("%T", decodedObject)))
			returnErrorResponse(w, "Failed to decode response data while checking voting permissions.")
			return
		}

		if err := proto.Unmarshal(resp.Data.GetValue(), any(decodedObject).(proto.Message)); err != nil {
			zap.L().Error("Failed to unmarshal protobuf", zap.Error(err))
			returnErrorResponse(w, "Failed to decode response data while checking voting permissions.")
			return
		}

		if decodedObject.Thread.UserId == r.Context().Value(middleware.UserIdKey).(string) {
			// do nothing
		} else if decodedObject.Thread.ReplyPermission == "Accounts You Follow" {
			userConn := getUserServiceConn()
			userClient := userpb.NewUserServiceClient(userConn)

			ctx, cancel := createContext()
			defer cancel()

			userResp, err := userClient.User_IsUserFollowing(ctx, &userpb.IsUserFollowingRequest{
				FollowingId: r.Context().Value(middleware.UserIdKey).(string),
				FollowerId:  decodedObject.Thread.UserId,
			})

			if err != nil || !userResp.Value {
				returnErrorResponse(w, "You are not allowed to vote on this thread. You must be followed by the author of this thread to reply.")
				return
			}

		} else if decodedObject.Thread.ReplyPermission == "Verified Accounts" {
			userConn := getUserServiceConn()
			userClient := userpb.NewUserServiceClient(userConn)

			ctx, cancel := createContext()
			defer cancel()

			userResp, _ := userClient.User_GetUserId(ctx, &userpb.GetProfileRequest{
				UserId: r.Context().Value(middleware.UserIdKey).(string),
			})
			decodedObject := new(userpb.UserProfile)
			if _, ok := any(decodedObject).(proto.Message); !ok {
				zap.L().Error("Type does not implement proto.Message", zap.String("type", fmt.Sprintf("%T", decodedObject)))
				returnErrorResponse(w, "Failed to decode response data while checking reply permissions.")
				return
			}

			if err := proto.Unmarshal(userResp.Data.GetValue(), any(decodedObject).(proto.Message)); err != nil {
				zap.L().Error("Failed to unmarshal protobuf", zap.Error(err))
				returnErrorResponse(w, "Failed to decode response data while checking reply permissions.")
				return
			}
			zap.L().Info("Decoded user profile", zap.String("userId", decodedObject.Username), zap.Bool("isVerified", decodedObject.IsVerified))
			if !decodedObject.IsVerified {
				returnErrorResponse(w, "You are not allowed vote on to this thread. You must be a verified account to reply.")
				return
			}
		}
	}

	ctx, cancel = createContext()
	defer cancel()

	resp, err = client.Thread_VoteThread(ctx, req)
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

func Thread_GetBookmarkedThreads(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (GetBookmarkedThreads) is called.")

	req, client := processThreadRequest[pb.StringThread](r, w)
	req.Value = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_GetBookmarkedThreads(ctx, req)
	processThreadResponseWithPayload[pb.GetThreadsResponse](resp, err, w)
}

func Thread_GetRepostedThreads(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (GetRepostedThreads) is called.")

	req, client := processThreadRequest[pb.StringThread](r, w)
	req.Value = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_GetRepostedThreads(ctx, req)
	processThreadResponseWithPayload[pb.GetThreadsResponse](resp, err, w)
}

func Thread_GetUserThreads(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (GetUserThreads) is called.")

	req, client := processThreadRequest[pb.UserToUserRequest](r, w)
	req.RequesterId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_GetUserThreads(ctx, req)
	processThreadResponseWithPayload[pb.GetThreadsResponse](resp, err, w)
}

func Thread_GetUserLikedThreads(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (GetUserLikedThreads) is called.")

	req, client := processThreadRequest[pb.StringThread](r, w)
	req.Value = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_GetUserLikedThreads(ctx, req)
	processThreadResponseWithPayload[pb.GetThreadsResponse](resp, err, w)
}

func Thread_GetUserReplies(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (GetUserReplies) is called.")

	req, client := processThreadRequest[pb.UserToUserRequest](r, w)
	req.RequesterId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_GetUserReplies(ctx, req)
	processThreadResponseWithPayload[pb.GetThreadsResponse](resp, err, w)
}

func Thread_GetUserMediaThreads(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (GetUserMediaThreads) is called.")

	req, client := processThreadRequest[pb.UserToUserRequest](r, w)
	req.RequesterId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_GetUserMediaThreads(ctx, req)
	processThreadResponseWithPayload[pb.GetThreadsResponse](resp, err, w)
}

func Thread_GetCommunityMediaThreads(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (GetCommunityMediaThreads) is called.")

	req, client := processThreadRequest[pb.GeneralThreadRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_GetCommunityMediaThreads(ctx, req)
	processThreadResponseWithPayload[pb.GetThreadsResponse](resp, err, w)
}

func Thread_GetTrendingHashtags(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (GetTrendingHashtags) is called.")

	conn := getThreadServiceConn()
	client := pb.NewThreadServiceClient(conn)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_GetTrendingHashtags(ctx, &pb.StringThread{})
	processThreadResponseWithPayload[pb.GetTrendingHashtagsResponse](resp, err, w)
}

func Thread_GetCommunityThreads(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (GetCommunityThreads) is called.")

	req, client := processThreadRequest[pb.GeneralThreadRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_GetCommunityThreads(ctx, req)
	processThreadResponseWithPayload[pb.GetThreadsResponse](resp, err, w)
}

func Thread_GetAdvertisementThreads(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (GetAdvertisementThreads) is called.")

	req, client := processThreadRequest[pb.GeneralThreadRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_GetAdvertisementThreads(ctx, req)
	processThreadResponseWithPayload[pb.GetThreadsResponse](resp, err, w)
}

func Thread_GetThreadByhashtag(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (GetThreadByhashtag) is called.")

	req, client := processThreadRequest[pb.GeneralThreadRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_GetThreadsByHashtag(ctx, req)
	processThreadResponseWithPayload[pb.GetThreadsResponse](resp, err, w)
}
func Thread_SearchThread(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (GetThreadByhashtag) is called.")

	req, client := processThreadRequest[pb.GeneralThreadRequest](r, w)
	req.UserId = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_SearchThread(ctx, req)
	processThreadResponseWithPayload[pb.GetThreadsResponse](resp, err, w)
}

func Thread_GetThreadCategories(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (GetCategories) is called.")

	conn := getThreadServiceConn()
	client := pb.NewThreadServiceClient(conn)

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Thread_GetThreadCategories(ctx, &pb.StringThread{})

	if err != nil {
		zap.L().Error("Error fetching thread categories", zap.Error(err))
		returnErrorResponse(w, "Error fetching thread categories: "+err.Error())
		return
	}

	response := &types.ApiResponse{
		Success: true,
		Message: "Thread categories retrieved successfully.",
		Payload: resp.Categories,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
