package bot

import (
	"SnLbot/internal/config"
	"SnLbot/internal/pkg/utils"
	"fmt"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbotapi.BotAPI
	cfg *config.Config
	//db  *sql.DB
	r *Router
}

func NewBot(cfg *config.Config /*, db *sql.DB */) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}
	api.Debug = false
	utils.LogInfo("Authorized on account: @%s", api.Self.UserName)

	bot := &Bot{
		api: api,
		cfg: cfg,
		// db: db,
	}

	bot.r = NewRouter()
	bot.registerHandlers()

	return bot, nil
}

func (bot *Bot) Start() {
	mode := bot.cfg.Mode
	if mode == "local" {
		bot.startPolling()
		return
	}

	if bot.cfg.WebhookURL != "" {
		bot.startWebhook()
	} else {
		bot.startPolling()
	}
}

func (bot *Bot) processMessage(msg *tgbotapi.Message) {
	if h, ok := bot.r.Resolve(msg.Text); ok {
		h(bot, msg)
		return
	}

	switch GetState(msg.Chat.ID) {
	case StateFAQ:
		bot.api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Пожалуйста, используйте кнопки меню или нажмите 'назад'"))
		return
	case StateCourses:
		courseWIPHandler(bot, msg)
		return
	}

	startHandler(bot, msg)
}

func (bot *Bot) startPolling() {
	utils.LogInfo("Running in POLLING mode")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10
	updates := bot.api.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			bot.processMessage(update.Message)
		}
	}
}

func (bot *Bot) startWebhook() {
	utils.LogInfo("Running in WEBHOOK mode")

	webhookURL := fmt.Sprintf("%s/%s", bot.cfg.WebhookURL, bot.api.Token)

	_, _ = bot.api.Request(tgbotapi.DeleteWebhookConfig{})

	wh, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		log.Fatalf("Failed to build webhook: %v", err)
	}

	_, err = bot.api.Request(wh)
	if err != nil {
		log.Fatalf("Failed to set webhook: %v", err)
	}

	updates := bot.api.ListenForWebhook("/" + bot.api.Token)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	go func() {
		for update := range updates {
			if update.Message != nil {
				bot.processMessage(update.Message)
			}
		}
	}()

	utils.LogInfo("Starting HTTP server on port %s", bot.cfg.Port)
	log.Fatal(http.ListenAndServe(":"+bot.cfg.Port, nil))
}

func (bot *Bot) registerHandlers() {
	commandMap := map[string]Handler{
		"/start": startHandler,
		"/help":  startHandler,

		// main menu
		"Частые вопросы": menuHandler("faq", StateFAQ),
		"Все курсы":      menuHandler("courses", StateCourses),
		"назад":          startHandler,

		// faq
		"Как проходит обучение?": faqHowHandler,
		"Формат обучения":        faqFormatHandler,

		// courses
		"Фигура человека": courseWIPHandler,
		"Форма и тон":     courseWIPHandler,
		"Свет и цвет":     courseWIPHandler,
	}

	for k, h := range commandMap {
		bot.r.RegisterCommand(k, h)
	}
}
