package middleware

import (
	"context"
	"net/http"

	"github.com/suryansh74/chat_app/pkg/logger"
	"github.com/suryansh74/chat_app/shared/helper"
	"github.com/suryansh74/chat_app/shared/token"
)

func VerifiedMiddleware(tokenMaker token.Maker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Log.Info("VerifiedMiddleware: checking auth and verified", "url", r.URL.Path)

			cookie, err := r.Cookie("session_token")
			if err != nil {
				logger.Log.Warn("VerifiedMiddleware: no session cookie", "error", err.Error())
				helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "unauthorized: session cookie missing",
				})
				return
			}

			tokenStr := cookie.Value

			payload, err := tokenMaker.VerifyToken(tokenStr)
			if err != nil {
				logger.Log.Warn("VerifiedMiddleware: invalid token", "error", err.Error())
				helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "unauthorized: invalid session",
				})
				return
			}

			if !payload.User.IsVerified {
				logger.Log.Warn("VerifiedMiddleware: user not verified", "email", payload.User.Email)
				helper.WriteJSON(w, http.StatusForbidden, map[string]string{
					"error": "forbidden: email not verified",
				})
				return
			}

			logger.Log.Info("VerifiedMiddleware: user verified", "user_id", payload.User.ID, "email", payload.User.Email)

			ctx := context.WithValue(r.Context(), UserContextKey, payload)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
