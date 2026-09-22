package middleware

import (
	"fmt"
	"strings"
	"context"
	"net/http"
	"6-project/pkg/jwt"
	"6-project/configs"
	"6-project/pkg/resp"
	"6-project/pkg/contextKeys"
)

func IsAuthed(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			resp.Json(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			resp.Json(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		isValid, data := jwt.NewJWT(config.Auth.Secret).Parse(token)
		if !isValid {
			resp.Json(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		fmt.Println("token:", isValid)
		fmt.Println("user:", data.Phone)
		ctx := context.WithValue(r.Context(), contextkey.UserContextKey, data.Phone)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}