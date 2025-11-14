package main

import (
	"SnLbot/internal/bot"
	"SnLbot/internal/config"
	"SnLbot/internal/pkg/utils"
	"log"
)

func main() {
	cfg := config.Load()

	utils.InitLogger(cfg)

	// dbConn, err := db.Connect(cfg)
	// if err != nil {
	// log.Fatalf("DB connection error: %v", err)
	// }
	// defer dbConn.Close()

	tgBot, err := bot.NewBot(cfg /* , dbConn */)
	if err != nil {
		log.Fatalf("Bot init error: %v", err)
	}

	tgBot.Start()
}
