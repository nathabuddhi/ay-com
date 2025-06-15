package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/middleware"
	communitypb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/community"
	threadpb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/thread"
	userpb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/user"
	"go.uber.org/zap"
)

func Admin_GetAllVerifyAccountRequest(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Admin Fetch Verify Account Request is called.")
	req, client := processUserRequest[userpb.StringUser](r, w)
	req.Value = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.Admin_GetAllVerifyAccountRequest(ctx, req)
	processUserResponseWithPayload[userpb.GetAllVerifyAccountResponse](resp, err, w)
}

func Admin_RejectPremiumRequest(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Admin Reject Premium Request is called.")
	req, client := processUserRequest[userpb.RejectPremiumRequest](r, w)

	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.Admin_RejectPremiumRequest(ctx, req)
	processUserResponseWithPayload[userpb.ApiResponseUser](resp, err, w)
}

func Admin_ApprovePremiumRequest(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Admin Approve Premium Request is called.")
	req, client := processUserRequest[userpb.StringUser](r, w)

	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.Admin_ApprovePremiumRequest(ctx, req)
	processUserResponseWithPayload[userpb.ApiResponseUser](resp, err, w)
}

func Admin_DeleteThread(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Thread (AdminDeleteThread) is called.")

	vars := mux.Vars(r)
	threadId := vars["id"]

	conn := getThreadServiceConn()
	client := threadpb.NewThreadServiceClient(conn)

	req := &threadpb.StringThread{}
	req.Value = threadId

	ctx, cancel := createContext()
	defer cancel()

	resp, err := client.Admin_DeleteThread(ctx, req)
	processThreadResponseWithoutPayload(resp, err, w)
}

func Admin_DeleteThreadCategory(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Admin adding thread category.")
	req, client := processThreadRequest[threadpb.StringThread](r, w)

	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.Admin_DeleteCategory(ctx, req)
	processThreadResponseWithoutPayload(resp, err, w)
}

func Admin_AddThreadCategory(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Admin adding thread category.")
	req, client := processThreadRequest[threadpb.StringThread](r, w)

	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.Admin_AddCategory(ctx, req)
	processThreadResponseWithoutPayload(resp, err, w)
}

func Admin_GetAllUsers(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Admin Get All Users is called.")
	conn := getUserServiceConn()
	client := userpb.NewUserServiceClient(conn)

	req := &userpb.StringUser{}
	req.Value = r.Context().Value(middleware.UserIdKey).(string)

	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.Admin_GetAllUsers(ctx, req)
	processUserResponseWithPayload[userpb.AdminUserResponse](resp, err, w)
}

func Admin_ToggleUserBan(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Admin toggling user ban category.")
	req, client := processUserRequest[userpb.StringUser](r, w)

	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.Admin_ToggleUserBan(ctx, req)
	processUserResponseWithoutPayload(resp, err, w)
}

func Admin_SendNewsLetter(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Admin sending newsletter.")
	req, client := processUserRequest[userpb.SendNewsLetterRequest](r, w)

	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.Admin_SendNewsLetter(ctx, req)
	processUserResponseWithoutPayload(resp, err, w)
}

func Admin_GetAllReports(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Admin getting all reports.")
	req := &userpb.StringUser{}
	req.Value = r.Context().Value(middleware.UserIdKey).(string)
	conn := getUserServiceConn()
	client := userpb.NewUserServiceClient(conn)

	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.Admin_GetAllReports(ctx, req)
	processUserResponseWithPayload[userpb.GetAllReportsResponse](resp, err, w)
}

func Admin_ApproveReport(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Admin approving report.")
	req, client := processUserRequest[userpb.StringUser](r, w)

	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.Admin_ApproveReport(ctx, req)
	processUserResponseWithoutPayload(resp, err, w)
}

func Admin_RejectReport(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Admin rejeting report.")
	req, client := processUserRequest[userpb.StringUser](r, w)

	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.Admin_RejectReport(ctx, req)
	processUserResponseWithoutPayload(resp, err, w)
}

func Admin_GetAllCommunityRequests(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Admin getting all community requests.")

	req := &communitypb.StringCommunity{}
	conn := getCommunityServiceConn()
	client := communitypb.NewCommunityServiceClient(conn)

	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.Admin_GetAllCommunityRequests(ctx, req)
	processCommunityResponseWithPayload[communitypb.GetCommunitiesResponse](resp, err, w)
}

func Admin_RejectCommunity(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Admin rejecting community.")
	req, client := processCommunityRequest[communitypb.RejectCommunityRequest](r, w)

	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.Admin_RejectCommunity(ctx, req)
	processCommunityResponseWithoutPayload(resp, err, w)
}

func Admin_ApproveCommunity(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Admin approving community.")
	req, client := processCommunityRequest[communitypb.StringCommunity](r, w)

	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.Admin_ApproveCommunity(ctx, req)
	processCommunityResponseWithoutPayload(resp, err, w)
}

func Admin_DeleteCommunityCategory(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Admin adding community category.")
	req, client := processCommunityRequest[communitypb.StringCommunity](r, w)

	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.Admin_DeleteCategory(ctx, req)
	processCommunityResponseWithoutPayload(resp, err, w)
}

func Admin_AddCommunityCategory(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Admin adding community category.")
	req, client := processCommunityRequest[communitypb.StringCommunity](r, w)

	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.Admin_AddCategory(ctx, req)
	processCommunityResponseWithoutPayload(resp, err, w)
}
