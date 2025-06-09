package handlers

import (
	"github.com/gorilla/mux"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

func InitRoutes() (r *mux.Router) {
	r = mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	InitUserRoutes(r)
	InitPublicThreadRoutes(r)
	zap.L().Info("Public Routes Initialized.")

	InitSecuredRoutes(r)
	return r
}

func InitAdminRoutes(r *mux.Router) {

	admin := r.PathPrefix("/").Subrouter()

	admin.Use(IsUserAdmin)
	admin.HandleFunc("/admin/deletethread", Admin_DeleteThread).Methods("DELETE")
}

func InitSecuredRoutes(r *mux.Router) {
	secured := r.PathPrefix("/").Subrouter()
	secured.Use(middleware.JwtAuthMiddleware)
	InitSecuredUserRoutes(secured)
	InitSecuredNotificationRoutes(secured)
	InitSecuredThreadRoutes(secured)
	zap.L().Info("Secured Routes Initialized.")

	// last
	InitAdminRoutes(secured)
	zap.L().Info("Admin Routes Initialized.")
}

func InitUserRoutes(r *mux.Router) {
	r.HandleFunc("/user/login", User_Login).Methods("POST")
	r.HandleFunc("/user/register", User_Register).Methods("POST")
	r.HandleFunc("/user/requestverificationcode", User_RequestVerificationCode).Methods("POST")
	r.HandleFunc("/user/validateverificationcode", User_ValidateVerificationCode).Methods("POST")
	r.HandleFunc("/user/getsecurityquestion", User_GetSecurityQuestion).Methods("POST")
	r.HandleFunc("/user/validatesecurityanswer", User_ValidateSecurityAnswer).Methods("POST")
	r.HandleFunc("/user/resetpassword", User_ResetPassword).Methods("PUT")
	r.HandleFunc("/user/refreshtoken", User_RefreshToken).Methods("POST")
}

func InitSecuredUserRoutes(secured *mux.Router) {
	secured.HandleFunc("/user/getprofile/{username}", User_GetProfile).Methods("GET")
	secured.HandleFunc("/user/getuserid", User_GetUserId).Methods("POST")
	secured.HandleFunc("/user/getselfprofile", User_GetSelfProfile).Methods("POST")
	secured.HandleFunc("/user/searchpeople", User_SearchPeople).Methods("POST")
	secured.HandleFunc("/user/changepassword", User_ChangePassword).Methods("PATCH")
	secured.HandleFunc("/user/updateprofile", User_UpdateProfile).Methods("PATCH")

	secured.HandleFunc("/user/followuser", User_FollowUser).Methods("POST")
	secured.HandleFunc("/user/unfollowuser", User_UnFollowUser).Methods("POST")
	secured.HandleFunc("/user/blockuser", User_BlockUser).Methods("POST")
	secured.HandleFunc("/user/unblockuser", User_UnBlockUser).Methods("POST")
	secured.HandleFunc("/user/getsettings", User_GetSettings).Methods("POST")
	secured.HandleFunc("/user/updatesettings", User_UpdateSettings).Methods("PATCH")

	secured.HandleFunc("/user/submitverifyaccountrequest", User_SubmitVerifyAccountRequest).Methods("POST")
	secured.HandleFunc("/user/getallverifyaccountrequest", User_GetAllVerifyAccountRequest).Methods("POST")

	secured.HandleFunc("/user/getallfollowers/{id}", User_GetAllFollowers).Methods("GET")
	secured.HandleFunc("/user/getallfollowing/{id}", User_GetAllFollowing).Methods("GET")
	secured.HandleFunc("/user/getallblocked", User_GetAllBlocked).Methods("POST")

	secured.HandleFunc("/user/changeavatar", User_ChangeAvatar).Methods("POST")
	secured.HandleFunc("/user/changebanner", User_ChangeBanner).Methods("POST")

	secured.HandleFunc("/user/deactivateaccount", User_DeactivateAccount).Methods("POST")

	secured.HandleFunc("/user/getfollowrecommendations", User_GetFollowRecommendations).Methods("GET")
}

func InitSecuredNotificationRoutes(secured *mux.Router) {
	secured.HandleFunc("/notification/getallnotifications", Notification_GetAllNotifications).Methods("POST")
	secured.HandleFunc("/notification/clearnotifications", Notification_ClearNotifications).Methods("DELETE")
	secured.HandleFunc("/notification/deletenotification", Notification_DeleteNotification).Methods("DELETE")
	secured.HandleFunc("/notification/marknotificationasread", Notification_MarkNotificationAsRead).Methods("PATCH")

	secured.HandleFunc("/notification/getsettings", Notification_GetSettings).Methods("POST")
	secured.HandleFunc("/notification/updatesettings", Notification_UpdateSettings).Methods("PATCH")
}

func InitPublicThreadRoutes(r *mux.Router) {
	r.HandleFunc("/thread/gettrendinghashtags", Thread_GetTrendingHashtags).Methods("GET")
}

func InitSecuredThreadRoutes(secured *mux.Router) {
	secured.HandleFunc("/thread/getallthreads", Thread_GetAllThreads).Methods("POST")
	secured.HandleFunc("/thread/getfollowingthreads", Thread_GetFollowingThreads).Methods("POST")
	secured.HandleFunc("/thread/create", Thread_CreateThread).Methods("POST")
	secured.HandleFunc("/thread/get/{id}", Thread_GetThreadById).Methods("GET")

	secured.HandleFunc("/thread/search", Thread_SearchThreads).Methods("POST")
	secured.HandleFunc("/thread/delete", Thread_DeleteThread).Methods("DELETE")

	secured.HandleFunc("/thread/pinthread", Thread_TogglePinThread).Methods("POST")

	secured.HandleFunc("/thread/vote", Thread_VoteThread).Methods("POST")

	secured.HandleFunc("/thread/togglelike", Thread_ToggleLike).Methods("POST")
	secured.HandleFunc("/thread/togglebookmark", Thread_ToggleBookmark).Methods("POST")
	secured.HandleFunc("/thread/togglerepost", Thread_ToggleRepost).Methods("POST")

	secured.HandleFunc("/thread/getbookmarks", Thread_GetBookmarkedThreads).Methods("POST")
	secured.HandleFunc("/thread/getreposts", Thread_GetRepostedThreads).Methods("POST")

	secured.HandleFunc("/thread/getuserthreads", Thread_GetUserThreads).Methods("POST")
	secured.HandleFunc("/thread/getuserlikedthreads", Thread_GetUserLikedThreads).Methods("POST")
	secured.HandleFunc("/thread/getuserreplies", Thread_GetUserReplies).Methods("POST")
	secured.HandleFunc("/thread/getusermediathreads", Thread_GetUserMediaThreads).Methods("POST")

	// secured.HandleFunc("/thread/getthreadsbyuser", Thread_).Methods("POST")
	// secured.HandleFunc("/thread/getthreadsbycommunity", Thread_).Methods("POST")
	// secured.HandleFunc("/thread/getpopularhashtags", Thread_).Methods("POST")
	// secured.HandleFunc("/thread/getthreadsbyhashtag", Thread_).Methods("POST")
	// secured.HandleFunc("/thread/getthreadsbymedia", Thread_).Methods("POST")
}
