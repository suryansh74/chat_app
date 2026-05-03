package middleware

import (
	"context"
	"net/http"

	"github.com/suryansh74/chat_app/shared/helper"
	"github.com/suryansh74/chat_app/shared/token"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(tokenMaker token.Maker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_token")
			if err != nil {
				helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "unauthorized: session cookie missing",
				})
				return
			}

			tokenStr := cookie.Value

			payload, err := tokenMaker.VerifyToken(tokenStr)
			if err != nil {
				helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "unauthorized: invalid session",
				})
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, payload)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GuestMiddleware(tokenMaker token.Maker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_token")
			if err == nil && cookie.Value != "" {
				tokenStr := cookie.Value
				_, err := tokenMaker.VerifyToken(tokenStr)
				if err == nil {
					helper.WriteJSON(w, http.StatusForbidden, map[string]string{
						"error": "forbidden: already logged in",
					})
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
