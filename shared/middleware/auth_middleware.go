package middleware

import (
	"context"
	"net/http"

	"github.com/suryansh74/chat_app/pkg/logger"
	"github.com/suryansh74/chat_app/shared/helper"
	"github.com/suryansh74/chat_app/shared/token"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(tokenMaker token.Maker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Log.Info("AuthMiddleware: checking auth", "url", r.URL.Path)

			cookie, err := r.Cookie("session_token")
			if err != nil {
				logger.Log.Warn("AuthMiddleware: no session cookie", "error", err.Error())
				helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "unauthorized: session cookie missing",
				})
				return
			}

			tokenStr := cookie.Value
			logger.Log.Info("AuthMiddleware: token found", "token_length", len(tokenStr))

			payload, err := tokenMaker.VerifyToken(tokenStr)
			if err != nil {
				logger.Log.Warn("AuthMiddleware: invalid token", "error", err.Error())
				helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "unauthorized: invalid session",
				})
				return
			}

			logger.Log.Info("AuthMiddleware: token valid", "user_id", payload.User.ID, "email", payload.User.Email)

			ctx := context.WithValue(r.Context(), UserContextKey, payload)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GuestMiddleware(tokenMaker token.Maker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Log.Info("GuestMiddleware: checking guest", "url", r.URL.Path)

			cookie, err := r.Cookie("session_token")
			if err == nil && cookie.Value != "" {
				tokenStr := cookie.Value
				_, err := tokenMaker.VerifyToken(tokenStr)
				if err == nil {
					logger.Log.Info("GuestMiddleware: user already logged in")
					helper.WriteJSON(w, http.StatusForbidden, map[string]string{
						"error": "forbidden: already logged in",
					})
					return
				}
			}

			logger.Log.Info("GuestMiddleware: allowing guest access")
			next.ServeHTTP(w, r)
		})
	}
}
