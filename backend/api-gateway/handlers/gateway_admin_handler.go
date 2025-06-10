package handlers

import (
	"net/http"

	"github.com/nathabuddhi/ay-com/backend/api-gateway/middleware"
	pb "github.com/nathabuddhi/ay-com/backend/api-gateway/proto/user"
	"go.uber.org/zap"
)

func Admin_GetAllVerifyAccountRequest(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Admin Fetch Verify Account Request is called.")
	req, client := processUserRequest[pb.StringUser](r, w)
	req.Value = r.Context().Value(middleware.UserIdKey).(string)
	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.Admin_GetAllVerifyAccountRequest(ctx, req)
	processUserResponseWithPayload[pb.GetAllVerifyAccountResponse](resp, err, w)
}

func Admin_RejectPremiumRequest(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Admin Reject Premium Request is called.")
	req, client := processUserRequest[pb.RejectPremiumRequest](r, w)

	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.Admin_RejectPremiumRequest(ctx, req)
	processUserResponseWithPayload[pb.ApiResponseUser](resp, err, w)
}

func Admin_ApprovePremiumRequest(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("Admin Approve Premium Request is called.")
	req, client := processUserRequest[pb.StringUser](r, w)

	ctx, cancel := createContext()
	defer cancel()
	resp, err := client.Admin_ApprovePremiumRequest(ctx, req)
	processUserResponseWithPayload[pb.ApiResponseUser](resp, err, w)
}
