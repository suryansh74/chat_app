package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/suryansh74/chat_app/config"
	"github.com/suryansh74/chat_app/internal/auth/handlers"
	authrepositories "github.com/suryansh74/chat_app/internal/auth/repositories"
	authservices "github.com/suryansh74/chat_app/internal/auth/services"
	"github.com/suryansh74/chat_app/pkg/logger"
	"github.com/suryansh74/chat_app/shared/token"
)

type server struct {
	cfg         *config.Config
	router      *chi.Mux
	authHandler *handlers.AuthHandler
	tokenMaker  token.Maker
}

func NewServer(cfg *config.Config) *server {
	repo := authrepositories.NewInMemoryUserRepository()
	service := authservices.NewAuthService(repo)
	tokenMaker, _ := token.NewPasetoMaker(cfg.TokenSymmetricKey)
	authHandler := handlers.NewAuthHandler(service, tokenMaker, cfg.CookieMaxAge, cfg.CookieSameSite)

	return &server{
		cfg:         cfg,
		router:      chi.NewRouter(),
		authHandler: authHandler,
		tokenMaker:  tokenMaker,
	}
}

func (s *server) Run() error {
	s.setupRoutes()

	addr := s.cfg.Host + ":" + s.cfg.Port
	logger.Log.Info("Server starting", "address", addr)
	if err := http.ListenAndServe(addr, s.router); err != nil {
		logger.Log.Fatal("Server failed to start", "error", err)
		return err
	}

	return nil
}
