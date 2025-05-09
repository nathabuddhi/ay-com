package handlers

import (
	"github.com/gorilla/mux"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/middleware"
	"go.uber.org/zap"
)

func InitRoutes() (r *mux.Router) {
	r = mux.NewRouter()

	InitUserRoutes(r)
	zap.L().Info("Public Routes Initialized.")

	InitSecuredRoutes(r)
	return r
}

func InitSecuredRoutes(r *mux.Router) {
	secured := r.PathPrefix("/").Subrouter()
	secured.Use(middleware.JwtAuthMiddleware)
	InitSecuredUserRoutes(secured)
	InitSecuredNotificationRoutes(secured)
	InitSecuredThreadRoutes(secured)

	zap.L().Info("Secured Routes Initialized.")
}

func InitUserRoutes(r *mux.Router) {
	r.HandleFunc("/user/login", User_Login).Methods("POST")
	r.HandleFunc("/user/register", User_Register).Methods("POST")
	r.HandleFunc("/user/requestverificationcode", User_RequestVerificationCode).Methods("POST")
	r.HandleFunc("/user/validateverificationcode", User_ValidateVerificationCode).Methods("POST")
	r.HandleFunc("/user/getsecurityquestion", User_GetSecurityQuestion).Methods("POST")
	r.HandleFunc("/user/validatesecurityanswer", User_ValidateSecurityAnswer).Methods("POST")
	r.HandleFunc("/user/resetpassword", User_ResetPassword).Methods("PUT")

}

func InitSecuredUserRoutes(secured *mux.Router) {
	secured.HandleFunc("/user/checktoken", User_CheckToken).Methods("POST")

	secured.HandleFunc("/user/getprofile/{username}", User_GetProfile).Methods("GET")
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
}

func InitSecuredNotificationRoutes(secured *mux.Router) {
	secured.HandleFunc("/notification/getallnotifications", Notification_GetAllNotifications).Methods("POST")
	secured.HandleFunc("/notification/clearnotifications", Notification_ClearNotifications).Methods("DELETE")
	secured.HandleFunc("/notification/deletenotification", Notification_DeleteNotification).Methods("DELETE")
	secured.HandleFunc("/notification/marknotificationasread", Notification_MarkNotificationAsRead).Methods("PATCH")

	secured.HandleFunc("/notification/getsettings", Notification_GetSettings).Methods("POST")
	secured.HandleFunc("/notification/updatesettings", Notification_UpdateSettings).Methods("PATCH")
}

func InitSecuredThreadRoutes(secured *mux.Router) {
	secured.HandleFunc("/thread/getallthreads", Thread_GetAllThreads).Methods("POST")
	// secured.HandleFunc("/thread/getthreadsbyuser", Thread_).Methods("POST")
	// secured.HandleFunc("/thread/getthreadsbycommunity", Thread_).Methods("POST")
	// secured.HandleFunc("/thread/searchthread", Thread_).Methods("POST")
	// secured.HandleFunc("/thread/getpopularhashtags", Thread_).Methods("POST")
	// secured.HandleFunc("/thread/getthreadsbyhashtag", Thread_).Methods("POST")
	// secured.HandleFunc("/thread/getthreadsbymedia", Thread_).Methods("POST")

	// secured.HandleFunc("/thread/createthread", Thread_).Methods("POST")
	// secured.HandleFunc("/thread/deletethread", Thread_).Methods("DELETE")

	// secured.HandleFunc("/thread/likethread", Thread_).Methods("POST")
	// secured.HandleFunc("/thread/unlikethread", Thread_).Methods("POST")

	// secured.HandleFunc("/thread/replythread", Thread_).Methods("POST")
	// secured.HandleFunc("/thread/deletereply", Thread_).Methods("POST")
}
