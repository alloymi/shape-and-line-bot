package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration loaded from env
type Config struct {
	BotToken   string
	Mode       string // "local" or "prod" (optional)
	Port       string
	WebhookURL string // set on Railway, e.g. https://myapp.railway.app

	// Postgres envs (not used yet)
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// Load reads .env and environment variables, returns Config
func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		BotToken:   os.Getenv("BOT_TOKEN"),
		Mode:       os.Getenv("MODE"),
		Port:       os.Getenv("PORT"),
		WebhookURL: os.Getenv("WEBHOOK_URL"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
	}

	if cfg.BotToken == "" {
		log.Fatal("BOT_TOKEN must be provided in environment")
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	// If Railway provides its URL variable, prefer that as WebhookURL
	if cfg.WebhookURL == "" {
		if v := os.Getenv("RAILWAY_STATIC_URL"); v != "" {
			cfg.WebhookURL = v
		}
	}

	return cfg
}
