package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken   string
	Mode       string
	Port       string
	WebhookURL string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	GoogleServiceAccountJSON string
	GoogleSpreadsheetID      string
}

func Load() *Config {
	_ = godotenv.Load() // ignore error — Railway использует Variables

	return &Config{
		BotToken:   os.Getenv("BOT_TOKEN"),
		Mode:       os.Getenv("MODE"),
		Port:       os.Getenv("PORT"),
		WebhookURL: os.Getenv("WEBHOOK_URL"),

		DBHost:     os.Getenv("PGHOST"),
		DBPort:     os.Getenv("PGPORT"),
		DBUser:     os.Getenv("PGUSER"),
		DBPassword: os.Getenv("PGPASSWORD"),
		DBName:     os.Getenv("PGDATABASE"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),

		GoogleServiceAccountJSON: os.Getenv("GOOGLE_SERVICE_ACCOUNT_JSON"),
		GoogleSpreadsheetID:      os.Getenv("GOOGLE_SPREADSHEET_ID"),
	}
}

func getEnvWithFallback(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
