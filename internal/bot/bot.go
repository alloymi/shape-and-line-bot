package bot

import (
	"SnLbot/internal/config"
	"SnLbot/internal/pkg/utils"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbotapi.BotAPI
	cfg *config.Config
	db  *sql.DB
	r   *Router
}

func NewBot(cfg *config.Config, db *sql.DB) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}
	api.Debug = false
	utils.LogInfo("Authorized on account: @%s", api.Self.UserName)

	bot := &Bot{
		api: api,
		cfg: cfg,
		db:  db,
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

	switch GetState(msg.Chat.ID) {
	case StateWaitlistChooseCourse:
		waitlistChooseCourseHandler(bot, msg)
		return
	case StateWaitlistAskFullName:
		waitlistFullNameHandler(bot, msg)
		return
	case StateWaitlistAskEmail:
		waitlistEmailHandler(bot, msg)
		return
	}

	if GetState(msg.Chat.ID) == StateCourseMenu {
		chatID := msg.Chat.ID
		course := userTempCourse[chatID]
		info := CoursesInfo[course]

		switch msg.Text {

		case "Длительность курса":
			bot.api.Send(tgbotapi.NewMessage(chatID, "Длительность: "+info.Duration))
			return

		case "Ближайший старт":
			bot.api.Send(tgbotapi.NewMessage(chatID, "Ближайший старт: "+info.StartDate))
			return

		case "Куратор курса":
			bot.api.Send(tgbotapi.NewMessage(chatID, info.Curator))
			return

		case "Записаться в лист ожидания":
			startWaitlistHandler(bot, msg)
			return

		case "Назад":
			ResetState(chatID)
			bot.api.Send(tgbotapi.NewMessage(chatID, "Выберите курс:"))
			back := tgbotapi.NewMessage(chatID, "")
			back.ReplyMarkup = Menus["courses"]
			bot.api.Send(back)
			return
		}
	}

	if h, ok := bot.r.Resolve(msg.Text); ok {
		h(bot, msg)
		return
	}

	switch GetState(msg.Chat.ID) {
	case StateFAQ:
		bot.api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Пожалуйста, используйте кнопки меню или нажмите 'назад'"))
		return
		//case StateCourses:
		//	courseWIPHandler(bot, msg)
		//	return
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
		"Частые вопросы":       menuHandler("faq", StateFAQ),
		"Все курсы":            menuHandler("courses", StateCourses),
		"Назад в главное меню": startHandler,

		// faq
		"Как проходит обучение?":                    faqHowHandler,
		"Формат обучения?":                          faqFormatHandler,
		"Хочу оплатить в рассрочку. Какие условия?": faqInstallmentHandler,
		"Я из другой страны. Могу ли я записаться на курс? Как проходит оплата?": faqForeignHandler,

		// courses
		"Фигура человека":                courseDetailsHandler,
		"Форма и тон":                    courseDetailsHandler,
		"Дизайн существ":                 courseDetailsHandler,
		"Портрет: Скетчинг и стилизация": courseDetailsHandler,
		"Свет и цвет":                    courseDetailsHandler,
		"Динамический портрет":           courseDetailsHandler,
		"Основы рисунка":                 courseDetailsHandler,
		"Мастерская с Евой":              courseDetailsHandler,
		"Анатомия человека":              courseDetailsHandler,

		//courses details
		"Длительность курса":    courseDurationHandler,
		"Ближайший старт":       courseStartHandler,
		"Куратор курса":         courseTeacherHandler,
		"Назад к списку курсов": courseBackHandler,

		// waiting list
		"Записаться в лист ожидания": startWaitlistHandler,
	}

	for k, h := range commandMap {
		bot.r.RegisterCommand(k, h)
	}
}

func isCourseName(s string) bool {
	courses := []string{
		"Фигура человека",
		"Форма и тон",
		"Дизайн существ",
		"Портрет: Скетчинг и стилизация",
		"Свет и цвет",
		"Динамический портрет",
		"Основы рисунка",
		"Мастерская с Евой",
		"Анатомия человека",
	}
	for _, c := range courses {
		if s == c {
			return true
		}
	}
	return false
}

func resetToMainMenu(b *Bot, chatID int64) {
	ResetState(chatID)
	delete(userTempCourse, chatID)
	delete(userTempFullname, chatID)

	msg := tgbotapi.NewMessage(chatID, "Запись отменена. Возвращение в главное меню.")
	msg.ReplyMarkup = Menus["main"]
	b.api.Send(msg)
}
