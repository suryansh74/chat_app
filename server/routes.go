package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/suryansh74/chat_app/internal/ws"
	"github.com/suryansh74/chat_app/pkg/logger"
	"github.com/suryansh74/chat_app/shared/helper"
	authmiddleware "github.com/suryansh74/chat_app/shared/middleware"
)

func (s *server) setupRoutes() {
	logger.Log.Info("Setting up CORS middleware...")
	s.router.Use(s.CORS)

	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Route("/api", func(r chi.Router) {
		r.Get("/check_health", checkHealth)

		r.Group(func(guest chi.Router) {
			guest.Use(authmiddleware.GuestMiddleware(s.tokenMaker))
			guest.Post("/auth/register", s.authHandler.Register)
			guest.Post("/auth/login", s.authHandler.Login)
		})

		r.Group(func(protected chi.Router) {
			protected.Use(authmiddleware.AuthMiddleware(s.tokenMaker))
			protected.Post("/auth/logout", s.authHandler.Logout)
		})

		r.Group(func(protected chi.Router) {
			protected.Use(authmiddleware.AuthMiddleware(s.tokenMaker))
			protected.Get("/profile", s.authHandler.Profile)
		})

		r.Group(func(emailVerified chi.Router) {
			emailVerified.Use(authmiddleware.AuthMiddleware(s.tokenMaker))
			emailVerified.Route("/email_verification", func(r chi.Router) {
				r.Post("/send_otp", s.emailVerificationHandler.SendOTP)
				r.Post("/verify_otp", s.emailVerificationHandler.VerifyOTP)
				r.Get("/verified", s.emailVerificationHandler.Verified)
			})
		})

		r.Group(func(passwordReset chi.Router) {
			passwordReset.Route("/password_reset", func(r chi.Router) {
				r.Post("/send_otp", s.passwordResetHandler.SendOTP)
				r.Post("/verify_otp", s.passwordResetHandler.VerifyOTP)
				r.Post("/set_password", s.passwordResetHandler.SetPassword)
			})
		})

		// Friends routes (requires verified email)
		r.Group(func(friends chi.Router) {
			friends.Use(authmiddleware.VerifiedMiddleware(s.tokenMaker))
			friends.Route("/friends", func(r chi.Router) {
				r.Get("/list", s.friendsHandler.GetFriends)
				r.Post("/request", s.friendsHandler.SendFriendRequest)
				r.Post("/accept", s.friendsHandler.AcceptFriendRequest)
				r.Post("/reject", s.friendsHandler.RejectFriendRequest)
				r.Get("/search", s.friendsHandler.SearchFriends)
				r.Delete("/", s.friendsHandler.RemoveFriend)
			})
		})

		// Search routes (requires verified email)
		r.Group(func(search chi.Router) {
			search.Use(authmiddleware.VerifiedMiddleware(s.tokenMaker))
			search.Route("/search", func(r chi.Router) {
				r.Get("/email", s.friendsHandler.SearchByEmail)
				r.Get("/global", s.chatHandler.SearchMessages)
				r.Get("/local", s.chatHandler.SearchConversationMessages)
			})
		})

		// Chat routes (requires verified email)
		r.Group(func(chat chi.Router) {
			chat.Use(authmiddleware.VerifiedMiddleware(s.tokenMaker))
			chat.Route("/chat", func(r chi.Router) {
				r.Get("/messages", s.chatHandler.GetMessages)
				r.Post("/messages", s.chatHandler.SendMessage)
			})
		})

		// Notification routes (requires verified email)
		r.Group(func(notification chi.Router) {
			notification.Use(authmiddleware.VerifiedMiddleware(s.tokenMaker))
			notification.Route("/notification", func(r chi.Router) {
				r.Get("/list", s.notificationHandler.GetNotifications)
				r.Get("/unread-count", s.notificationHandler.GetUnreadCount)
				r.Put("/read", s.notificationHandler.MarkAsRead)
				r.Put("/read-all", s.notificationHandler.MarkAllAsRead)
			})
		})
	})

	// WebSocket route (no auth middleware - auth via query param)
	s.router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(s.wsHub, w, r)
	})
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Check health")
	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "server is healthy",
	})
}
