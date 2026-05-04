package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	})
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Check health")
	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "server is healthy",
	})
}
