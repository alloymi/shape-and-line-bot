package utils

import (
	"SnLbot/internal/config"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger(cfg *config.Config) {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if cfg.Mode == "local" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func LogInfo(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}
func LogError(format string, v ...interface{}) {
	log.Error().Msgf(format, v...)
}
