package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/suryansh74/chat_app/config"
	"github.com/suryansh74/chat_app/internal/auth/handlers"
	authrepositories "github.com/suryansh74/chat_app/internal/auth/repositories"
	authservices "github.com/suryansh74/chat_app/internal/auth/services"
	chathandlers "github.com/suryansh74/chat_app/internal/chat/handlers"
	chatrepositories "github.com/suryansh74/chat_app/internal/chat/repositories"
	chatservices "github.com/suryansh74/chat_app/internal/chat/services"
	friendshandlers "github.com/suryansh74/chat_app/internal/friends/handlers"
	friendsrepositories "github.com/suryansh74/chat_app/internal/friends/repositories"
	friendsservices "github.com/suryansh74/chat_app/internal/friends/services"
	notificationhandlers "github.com/suryansh74/chat_app/internal/notification/handlers"
	notificationrepositories "github.com/suryansh74/chat_app/internal/notification/repositories"
	notificationservices "github.com/suryansh74/chat_app/internal/notification/services"
	ws "github.com/suryansh74/chat_app/internal/ws"
	"github.com/suryansh74/chat_app/pkg/logger"
	presencehandlers "github.com/suryansh74/chat_app/server/handlers"
	"github.com/suryansh74/chat_app/shared/cache"
	emailadapters "github.com/suryansh74/chat_app/shared/email/adapters"
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
	friendsHandler           *friendshandlers.FriendsHandler
	chatHandler              *chathandlers.ChatHandler
	notificationHandler      *notificationhandlers.NotificationHandler
	presenceHandler          *presencehandlers.PresenceHandler
	wsHub                    *ws.Hub
	tokenMaker               token.Maker
	redisClient              *redis.Client
}

func NewServer(cfg *config.Config) *server {
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

	if err := authrepositories.AutoMigrate(db); err != nil {
		logger.Log.Fatal("Failed to migrate database", "error", err)
	}
	logger.Log.Info("Database migrated successfully")

	if err := friendsrepositories.AutoMigrate(db); err != nil {
		logger.Log.Fatal("Failed to migrate friends", "error", err)
	}

	if err := chatrepositories.AutoMigrate(db); err != nil {
		logger.Log.Fatal("Failed to migrate chat", "error", err)
	}

	if err := notificationrepositories.AutoMigrate(db); err != nil {
		logger.Log.Fatal("Failed to migrate notifications", "error", err)
	}

	repo := authrepositories.NewMySQLUserRepository(db)
	service := authservices.NewAuthService(repo)
	tokenMaker, _ := token.NewPasetoMaker(cfg.TokenSymmetricKey)
	authHandler := handlers.NewAuthHandler(service, tokenMaker, cfg.CookieMaxAge, cfg.CookieSameSite)

	emailVerificationService := authservices.NewEmailVerificationService(repo, cfg.OtpExpiryMinutes, cfg.OtpMaxAttempts)

	var emailSender emailadapters.EmailSenderPort
	if cfg.SMTPUsername == "" && cfg.SMTPPassword == "" {
		logger.Log.Info("Using Mailpit adapter", "host", cfg.SMTPHost, "port", cfg.SMTPPort)
		emailSender = emailadapters.NewMailpitAdapter(cfg.SMTPHost, cfg.SMTPPort)
	} else {
		logger.Log.Info("Using Gmail adapter", "host", cfg.SMTPHost, "port", cfg.SMTPPort)
		emailSender = emailadapters.NewGmailAdapter(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword)
	}

	emailVerificationHandler := handlers.NewEmailVerificationHandler(emailVerificationService, emailSender, tokenMaker, cfg.CookieMaxAge)
	passwordResetHandler := handlers.NewPasswordResetHandler(emailVerificationService, emailSender, tokenMaker, cfg.PasswordResetRedirectURL)

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Log.Fatal("Failed to connect to Redis", "error", err)
	}
	logger.Log.Info("Connected to Redis", "url", cfg.RedisURL)

	redisCache := cache.NewRedisCacheFromClient(redisClient)

	friendRepo := friendsrepositories.NewMySQLFriendRepository(db)
	notificationRepo := notificationrepositories.NewMySQLNotificationRepository(db)
	notificationSvc := notificationservices.NewNotificationService(notificationRepo, redisCache, emailSender)

	wsHub := ws.NewHub(redisClient)
	go wsHub.Run()

	friendsService := friendsservices.NewFriendsService(friendRepo, notificationSvc, emailSender, tokenMaker, wsHub)
	friendsHandler := friendshandlers.NewFriendsHandler(friendsService)

	chatRepo := chatrepositories.NewMySQLChatRepository(db)
	chatService := chatservices.NewChatService(chatRepo, friendRepo, notificationSvc, wsHub)
	chatHandler := chathandlers.NewChatHandler(chatService)

	notificationHandler := notificationhandlers.NewNotificationHandler(notificationSvc)
	presenceHandler := presencehandlers.NewPresenceHandler(redisClient)

	return &server{
		cfg:                      cfg,
		router:                   chi.NewRouter(),
		authHandler:              authHandler,
		emailVerificationHandler: emailVerificationHandler,
		passwordResetHandler:     passwordResetHandler,
		friendsHandler:           friendsHandler,
		chatHandler:              chatHandler,
		notificationHandler:      notificationHandler,
		presenceHandler:          presenceHandler,
		wsHub:                    wsHub,
		tokenMaker:               tokenMaker,
		redisClient:              redisClient,
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
