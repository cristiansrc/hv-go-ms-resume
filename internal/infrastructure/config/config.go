package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Port    string
	Version string

	JWTSecret          string
	JWTExpirationHours int

	AdminUsername     string
	AdminPasswordHash string

	DatabasePath string

	AWSRegion    string
	AWSAccessKey string
	AWSSecretKey string
	AWSBucket    string

	RenderCVURL string

	TelegramBotToken string
	TelegramChatID   string

	LogLevel string

	RateLimitPublic        int
	RateLimitLogin         int
	RateLimitAuthenticated int
}

// Load reads configuration from environment variables with defaults.
func Load() (*Config, error) {
	jwtExpHours, err := getEnvInt("JWT_EXPIRATION_HOURS", 24)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRATION_HOURS: %w", err)
	}

	rateLimitPublic, err := getEnvInt("RATE_LIMIT_PUBLIC", 100)
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_PUBLIC: %w", err)
	}

	rateLimitLogin, err := getEnvInt("RATE_LIMIT_LOGIN", 5)
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_LOGIN: %w", err)
	}

	rateLimitAuth, err := getEnvInt("RATE_LIMIT_AUTHENTICATED", 200)
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_AUTHENTICATED: %w", err)
	}

	port := getEnv("PORT", "8080")
	version := getEnv("APP_VERSION", "1.0.0")
	databasePath := getEnv("DATABASE_PATH", "./data/resume.db")
	logLevel := getEnv("LOG_LEVEL", "info")

	jwtSecret, err := requireEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}
	adminUsername, err := requireEnv("ADMIN_USERNAME")
	if err != nil {
		return nil, err
	}
	adminPasswordHash, err := requireEnv("ADMIN_PASSWORD_HASH")
	if err != nil {
		return nil, err
	}
	awsRegion, err := requireEnv("AWS_REGION")
	if err != nil {
		return nil, err
	}
	awsAccessKey, err := requireEnv("AWS_ACCESS_KEY_ID")
	if err != nil {
		return nil, err
	}
	awsSecretKey, err := requireEnv("AWS_SECRET_ACCESS_KEY")
	if err != nil {
		return nil, err
	}
	awsBucket, err := requireEnv("AWS_BUCKET")
	if err != nil {
		return nil, err
	}
	renderCVURL, err := requireEnv("RENDER_CV_URL")
	if err != nil {
		return nil, err
	}
	telegramBotToken, err := requireEnv("TELEGRAM_BOT_TOKEN")
	if err != nil {
		return nil, err
	}
	telegramChatID, err := requireEnv("TELEGRAM_CHAT_ID")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Port:    port,
		Version: version,

		JWTSecret:          jwtSecret,
		JWTExpirationHours: jwtExpHours,

		AdminUsername:     adminUsername,
		AdminPasswordHash: adminPasswordHash,

		DatabasePath: databasePath,

		AWSRegion:    awsRegion,
		AWSAccessKey: awsAccessKey,
		AWSSecretKey: awsSecretKey,
		AWSBucket:    awsBucket,

		RenderCVURL: renderCVURL,

		TelegramBotToken: telegramBotToken,
		TelegramChatID:   telegramChatID,

		LogLevel: logLevel,

		RateLimitPublic:        rateLimitPublic,
		RateLimitLogin:         rateLimitLogin,
		RateLimitAuthenticated: rateLimitAuth,
	}

	return cfg, nil
}

// JWTExpiration returns the JWT expiration duration.
func (c *Config) JWTExpiration() time.Duration {
	return time.Duration(c.JWTExpirationHours) * time.Hour
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

func requireEnv(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return "", fmt.Errorf("required environment variable %s is not set", key)
	}
	return val, nil
}

func getEnvInt(key string, fallback int) (int, error) {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("cannot parse %s=%q: %w", key, val, err)
	}
	return n, nil
}
