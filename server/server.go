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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	// Connect to MySQL
	dsn := authrepositories.GetDSN(
		cfg.MySQLHost,
		cfg.MySQLPort,
		cfg.MySQLUser,
		cfg.MySQLPassword,
		cfg.MySQLDatabase,
	)

	logger.Log.Info("Connecting to MySQL...", "dsn", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal("Failed to connect to MySQL", "error", err)
	}

	// Auto migrate the database
	if err := authrepositories.AutoMigrate(db); err != nil {
		logger.Log.Fatal("Failed to migrate database", "error", err)
	}
	logger.Log.Info("Database migrated successfully")

	repo := authrepositories.NewMySQLUserRepository(db)
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

func (s *server) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Info("CORS middleware", "method", r.Method, "url", r.URL.Path, "origin", r.Header.Get("Origin"))

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Cookie")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			logger.Log.Info("CORS preflight handled", "url", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *server) Run() error {
	logger.Log.Info("Setting up routes...")
	s.setupRoutes()

	addr := s.cfg.Host + ":" + s.cfg.Port
	logger.Log.Info("Server starting", "address", addr)
	if err := http.ListenAndServe(addr, s.router); err != nil {
		logger.Log.Fatal("Server failed to start", "error", err)
		return err
	}

	return nil
}
