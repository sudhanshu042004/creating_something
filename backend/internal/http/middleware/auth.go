package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sudhanshu042004/sandbox/internal/utils/response"
	"github.com/sudhanshu042004/sandbox/internal/utils/token"
)

type contextKey string

const userIdKey = contextKey("userId")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("session")
		if tokenString == "" {
			response.WriteJson(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
			return
		}
		userId, err := token.VerifyToken(tokenString)
		if err != nil {
			response.WriteJson(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
			return
		}

		ctx := context.WithValue(r.Context(), userIdKey, userId)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func UserKeyId() contextKey {
	return userIdKey
}
