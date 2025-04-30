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

func InitUserRoutes(r *mux.Router) {
	r.HandleFunc("/user/login", User_Login).Methods("POST")
	r.HandleFunc("/user/register", User_Register).Methods("POST")
	r.HandleFunc("/user/requestverificationcode", User_RequestVerificationCode).Methods("POST")
	r.HandleFunc("/user/validateverificationcode", User_ValidateVerificationCode).Methods("POST")
}

func InitSecuredRoutes(r *mux.Router) {
	secured := r.PathPrefix("/").Subrouter()
	secured.Use(middleware.JwtAuthMiddleware)
	secured.HandleFunc("/user/getprofile/{id}", User_GetProfile).Methods("GET")
	secured.HandleFunc("/user/changepassword", User_ChangePassword).Methods("POST")
}
