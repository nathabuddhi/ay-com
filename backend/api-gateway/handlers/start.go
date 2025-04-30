package handlers

import (
	"github.com/gorilla/mux"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/middleware"
)

func InitRoutes() (r *mux.Router) {
	r = mux.NewRouter()

	InitUserRoutes(r)

	InitSecuredRoutes(r)
	return r
}

func InitSecuredRoutes(r *mux.Router) {
	secured := r.PathPrefix("/").Subrouter()
	secured.Use(middleware.JwtAuthMiddleware)
	InitSecuredUserRoutes(secured)
}

func InitUserRoutes(r *mux.Router) {
	r.HandleFunc("/user/login", User_Login).Methods("POST")
	r.HandleFunc("/user/register", User_Register).Methods("POST")
	r.HandleFunc("/user/requestverificationcode", User_RequestVerificationCode).Methods("POST")
	r.HandleFunc("/user/validateverificationcode", User_ValidateVerificationCode).Methods("POST")
	r.HandleFunc("/user/getsecurityquestion", User_GetSecurityQuestion).Methods("POST")
	r.HandleFunc("/user/validatesecurityanswer", User_ValidateSecurityAnswer).Methods("POST")
	r.HandleFunc("/user/resetpassword", User_ResetPassword).Methods("POST")

}

func InitSecuredUserRoutes(secured *mux.Router) {
	secured.HandleFunc("/user/getprofile/{id}", User_GetProfile).Methods("GET")
	secured.HandleFunc("/user/changepassword", User_ChangePassword).Methods("POST")
}
