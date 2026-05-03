package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Host string `mapstructure:"SERVER_HOST"`
	Port string `mapstructure:"SERVER_PORT"`

	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	CookieMaxAge        int           `mapstructure:"COOKIE_MAX_AGE"`
	CookieSameSite      string        `mapstructure:"COOKIE_SAMESITE"`

	// DBPort         string `mapstructure:"DB_PORT"`
	// DBHost         string `mapstructure:"DB_HOST"`
	// DBUser         string `mapstructure:"DB_USER"`
	// DBName         string `mapstructure:"DB_NAME"`
	// DBPassword     string `mapstructure:"DB_PASSWORD"`
	// CollectionName string `mapstructure:"COLLECTION_NAME"`
	//
	// GoogleClientID      string `mapstructure:"GOOGLE_CLIENT_ID"`
	// GoogleClientSecret  string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	// GoogleRedirectURL   string `mapstructure:"GOOGLE_REDIRECT_URL"`
	// FrontendURL         string `mapstructure:"FRONTEND_URL"`
	//
	// SMTPHost     string `mapstructure:"SMTP_HOST"`
	// SMTPPort     string `mapstructure:"SMTP_PORT"`
	// SMTPUsername string `mapstructure:"SMTP_USERNAME"`
	// SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
	//
	// OtpExpiryMinutes int `mapstructure:"OTP_EXPIRY_MINUTES"`
	// OtpMaxAttempts   int `mapstructure:"OTP_MAX_ATTEMPTS"`
	//
	// LogLevel string `mapstructure:"LOG_LEVEL"`
	// LogPath  string `mapstructure:"LOG_PATH"`
	//
	// EnableSwaggerUI bool `mapstructure:"ENABLE_SWAGGER_UI"`
}

func LoadConfig() (config Config, err error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("failed to read config: %w", err)
	}

	err = viper.Unmarshal(&config)
	return config, err
}
