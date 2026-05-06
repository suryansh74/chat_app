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

	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     string `mapstructure:"SMTP_PORT"`
	SMTPUsername string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`

	OtpExpiryMinutes         int    `mapstructure:"OTP_EXPIRY_MINUTES"`
	OtpMaxAttempts           int    `mapstructure:"OTP_MAX_ATTEMPTS"`
	PasswordResetRedirectURL string `mapstructure:"PASSWORD_RESET_REDIRECT_URL"`

	MySQLHost     string `mapstructure:"MYSQL_HOST"`
	MySQLPort     string `mapstructure:"MYSQL_PORT"`
	MySQLUser     string `mapstructure:"MYSQL_USER"`
	MySQLPassword string `mapstructure:"MYSQL_PASSWORD"`
	MySQLDatabase string `mapstructure:"MYSQL_DATABASE"`

	RedisURL string `mapstructure:"REDIS_URL"`
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
