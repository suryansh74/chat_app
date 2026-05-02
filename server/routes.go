package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/suryansh74/chat_app/pkg/logger"
	"github.com/suryansh74/chat_app/shared/helper"
)

func (s *server) setupRoutes() {
	s.router.Route("/api", func(r chi.Router) {
		r.Get("/check_health", checkHealth)
	})
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Check health")
	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "server is healthy",
	})
}
