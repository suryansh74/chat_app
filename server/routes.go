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
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Route("/api", func(r chi.Router) {
		r.Get("/check_health", checkHealth)

		r.Route("/auth", func(r chi.Router) {
			r.Use(authmiddleware.GuestMiddleware(s.tokenMaker))
			r.Post("/register", s.authHandler.Register)
			r.Post("/login", s.authHandler.Login)
		})

		r.Group(func(protected chi.Router) {
			protected.Use(authmiddleware.AuthMiddleware(s.tokenMaker))
			protected.Get("/profile", s.authHandler.Profile)
			protected.Post("/logout", s.authHandler.Logout)
		})
	})
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Check health")
	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "server is healthy",
	})
}
