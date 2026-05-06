package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App         AppConfig
	DB          DBConfig
	Redis       RedisConfig
	RabbitMQ    RabbitMQConfig
	Email       EmailConfig
	RateLimit   RateLimitConfig
	Idempotency IdempotencyConfig
}

type AppConfig struct {
	Port string
	Env  string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host string
	Port string
}

type RabbitMQConfig struct {
	URL string
}

type EmailConfig struct {
	Provider      string
	FromEmail     string
	FromName      string
	ResendAPIKey  string
	ResendTimeout time.Duration
	ResendBaseURL string
	SESRegion     string
}

type RateLimitConfig struct {
	Enabled bool
	Limit   int
	Window  time.Duration
}

type IdempotencyConfig struct {
	Enabled bool
	TTL     time.Duration
}

func Load() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("warning: no .env file found, reading from system environment")
	}

	cfg := &Config{
		App: AppConfig{
			Port: getEnvOptional("APP_PORT", "3001"),
			Env:  getEnvOptional("APP_ENV", "development"),
		},
		DB: DBConfig{
			Host:     getEnvRequired("DB_HOST"),
			Port:     getEnvRequired("DB_PORT"),
			User:     getEnvRequired("DB_USER"),
			Password: getEnvRequired("DB_PASSWORD"),
			Name:     getEnvRequired("DB_NAME"),
			SSLMode:  getEnvOptional("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host: getEnvOptional("REDIS_HOST", ""),
			Port: getEnvOptional("REDIS_PORT", ""),
		},
		RateLimit: RateLimitConfig{
			Enabled: getEnvOptionalBool("RATE_LIMIT_ENABLED", true),
			Limit:   getEnvOptionalInt("RATE_LIMIT_MAX_REQUESTS", 5),
			Window:  getEnvOptionalDuration("RATE_LIMIT_WINDOW", time.Minute),
		},
		Idempotency: IdempotencyConfig{
			Enabled: getEnvOptionalBool("IDEMPOTENCY_ENABLED", true),
			TTL:     getEnvOptionalDuration("IDEMPOTENCY_TTL", 10*time.Minute),
		},
		RabbitMQ: RabbitMQConfig{
			URL: getEnvOptional("RABBITMQ_URL", ""),
		},
		Email: EmailConfig{
			Provider:      getEnvOptional("EMAIL_PROVIDER", "resend"),
			FromEmail:     getEnvOptional("EMAIL_FROM_EMAIL", ""),
			FromName:      getEnvOptional("EMAIL_FROM_NAME", ""),
			ResendAPIKey:  getEnvOptional("RESEND_API_KEY", ""),
			ResendBaseURL: getEnvOptional("RESEND_BASE_URL", ""),
			ResendTimeout: getEnvOptionalDuration("RESEND_TIMEOUT", 10*time.Second),
			SESRegion:     getEnvOptional("AWS_SES_REGION", "us-east-1"),
		},
	}

	return cfg
}

func getEnvRequired(key string) string {
	value, exists := os.LookupEnv(key)

	if value == "" || !exists {
		panic(fmt.Sprintf("The value of key:%s is not available to read", key))
	}

	return value
}

func getEnvOptional(key, fallback string) string {
	value, exists := os.LookupEnv(key)

	if value == "" && !exists {
		return fallback
	}

	return value
}

func getEnvOptionalInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("invalid integer value for %s: %q", key, value))
	}

	return parsed
}

func getEnvOptionalBool(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		panic(fmt.Sprintf("invalid boolean value for %s: %q", key, value))
	}

	return parsed
}

func getEnvOptionalDuration(key string, defaultValue time.Duration) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		panic(fmt.Sprintf("invalid duration value for %s: %q", key, value))
	}

	return parsed
}
