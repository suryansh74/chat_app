package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/suryansh74/chat_app/config"
	"github.com/suryansh74/chat_app/pkg/logger"
)

type server struct {
	cfg    *config.Config
	router *chi.Mux
}

func NewServer(cfg *config.Config) *server {
	return &server{
		cfg:    cfg,
		router: chi.NewRouter(),
	}
}

func (s *server) Run() error {
	// setup routes
	s.setupRoutes()

	// running server
	addr := s.cfg.Host + ":" + s.cfg.Port
	logger.Log.Info("Server starting", "address", addr)
	if err := http.ListenAndServe(addr, s.router); err != nil {
		logger.Log.Fatal("Server failed to start", "error", err)
		return err
	}

	return nil
}
