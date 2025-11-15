package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib"

	"SnLbot/internal/config"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort,
		cfg.DBUser, cfg.DBPassword,
		cfg.DBName, cfg.DBSSLMode,
	)

	log.Printf("Connecting to database: host=%s, dbname=%s", cfg.DBHost, cfg.DBName)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Printf("Successfully connected to database: %s", cfg.DBName)
	return db, nil
}

func RunMigrations(db *sql.DB) error {
	migrationsPath := "internal/db/migrations"

	files, err := filepath.Glob(filepath.Join(migrationsPath, "*.sql"))
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %v", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no migration files found in %s", migrationsPath)
	}

	for _, file := range files {
		log.Printf("Running migration: %s", file)

		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading migration file %s: %v", file, err)
		}

		_, err = db.Exec(string(sqlBytes))
		if err != nil {
			return fmt.Errorf("error executing migration %s: %v", file, err)
		}

		log.Printf("Migration completed: %s", file)
	}

	return nil
}
