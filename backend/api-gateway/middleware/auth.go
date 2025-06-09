package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nathabuddhi/ay-com/backend/api-gateway/types"
)

var UserIdKey = &contextKey{"user_id"}

type contextKey struct {
	name string
}

func JwtAuthMiddleware(next http.Handler) http.Handler {
	var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response := types.ApiResponse{
				Success: false,
				Message: "No Token Provided.",
				Payload: nil,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			response := types.ApiResponse{
				Success: false,
				Message: "Invalid Token Provided.",
				Payload: nil,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response := types.ApiResponse{
				Success: false,
				Message: "Invalid Token Claims.",
				Payload: nil,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		userID := claims["user_id"].(string)

		ctx := context.WithValue(r.Context(), UserIdKey, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
