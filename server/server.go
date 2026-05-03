package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/suryansh74/chat_app/config"
	"github.com/suryansh74/chat_app/internal/auth/handlers"
	authrepositories "github.com/suryansh74/chat_app/internal/auth/repositories"
	authservices "github.com/suryansh74/chat_app/internal/auth/services"
	"github.com/suryansh74/chat_app/pkg/logger"
	"github.com/suryansh74/chat_app/shared/email"
	"github.com/suryansh74/chat_app/shared/token"
)

type server struct {
	cfg                      *config.Config
	router                   *chi.Mux
	authHandler              *handlers.AuthHandler
	emailVerificationHandler *handlers.EmailVerificationHandler
	passwordResetHandler     *handlers.PasswordResetHandler
	tokenMaker               token.Maker
}

func NewServer(cfg *config.Config) *server {
	repo := authrepositories.NewInMemoryUserRepository()
	service := authservices.NewAuthService(repo)
	tokenMaker, _ := token.NewPasetoMaker(cfg.TokenSymmetricKey)
	authHandler := handlers.NewAuthHandler(service, tokenMaker, cfg.CookieMaxAge, cfg.CookieSameSite)

	emailVerificationService := authservices.NewEmailVerificationService(repo, cfg.OtpExpiryMinutes, cfg.OtpMaxAttempts)
	emailSender := email.NewSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword)
	emailVerificationHandler := handlers.NewEmailVerificationHandler(emailVerificationService, emailSender, tokenMaker, cfg.CookieMaxAge)
	passwordResetHandler := handlers.NewPasswordResetHandler(emailVerificationService, emailSender, tokenMaker, cfg.PasswordResetRedirectURL)

	return &server{
		cfg:                      cfg,
		router:                   chi.NewRouter(),
		authHandler:              authHandler,
		emailVerificationHandler: emailVerificationHandler,
		passwordResetHandler:     passwordResetHandler,
		tokenMaker:               tokenMaker,
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
