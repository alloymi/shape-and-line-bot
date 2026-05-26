package main

import (
	"SnLbot/internal/db"
	"log"

	"SnLbot/internal/bot"
	"SnLbot/internal/config"
	"SnLbot/internal/pkg/utils"
)

func main() {
	cfg := config.Load()
	utils.InitLogger(cfg)

	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	defer dbConn.Close()

	err = db.RunMigrations(dbConn)
	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	tgBot, err := bot.NewBot(cfg, dbConn)
	//tgBot, err := bot.NewBot(cfg, nil)
	if err != nil {
		log.Fatalf("Bot init error: %v", err)
	}

	tgBot.Start()
}
