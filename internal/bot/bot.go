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

		case "Основная информация":
			bot.api.Send(tgbotapi.NewMessage(chatID, info.MainInfo))
			return

		case "Доступные тарифы":
			bot.api.Send(tgbotapi.NewMessage(chatID, info.Tariffs))
			return

		case "Программа курса":
			bot.api.Send(tgbotapi.NewMessage(chatID, info.Schedule))
			return

		case "О чем курс":
			bot.api.Send(tgbotapi.NewMessage(chatID, info.About))
			return

		case "Что понадобится":
			bot.api.Send(tgbotapi.NewMessage(chatID, info.Tools))
			return

		case "Для кого подходит курс":
			bot.api.Send(tgbotapi.NewMessage(chatID, info.ForWhom))
			return

		case "Где посмотреть работы куратора и студентов?":
			bot.api.Send(tgbotapi.NewMessage(chatID, info.WhereToFindWorks))
			return

		case "Назад к выбору курса":
			ResetState(chatID)
			bot.api.Send(tgbotapi.NewMessage(chatID, "Выберите курс:"))
			back := tgbotapi.NewMessage(chatID, "")
			back.ReplyMarkup = Menus["courses"]
			bot.api.Send(back)
			return
		}
	}

	if GetState(msg.Chat.ID) == StateFAQPayment ||
		GetState(msg.Chat.ID) == StateFAQStudy ||
		GetState(msg.Chat.ID) == StateFAQAbout {

		if msg.Text == "назад" {
			SetState(msg.Chat.ID, StateFAQ)
			resp := tgbotapi.NewMessage(msg.Chat.ID, "Выберите категорию вопросов:")
			resp.ReplyMarkup = faqCategoriesMenu()
			bot.api.Send(resp)
			return
		}

		if h, ok := bot.r.Resolve(msg.Text); ok {
			h(bot, msg)
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

	wh, _ := bot.api.GetWebhookInfo()
	if wh.IsSet() {
		utils.LogInfo("Webhook detected, deleting it...")
		_, err := bot.api.Request(tgbotapi.DeleteWebhookConfig{DropPendingUpdates: true})
		if err != nil {
			utils.LogError("Failed to delete webhook: %v", err)
		} else {
			utils.LogInfo("Webhook deleted successfully")
		}
	}
	utils.LogInfo("Running in POLLING mode")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
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
		"Частые вопросы":       faqHandler,
		"Все курсы":            menuHandler("courses", StateCourses),
		"Назад в главное меню": startHandler,

		// faq
		"О школе":             faqCategoryHandler,
		"Вопросы об оплате":   faqCategoryHandler,
		"Вопросы об обучении": faqCategoryHandler,

		"Подробнее про школу":                                                    faqAboutHandler,
		"Как проходит обучение?":                                                 faqFormatHandler,
		"Как я могу записаться на курс?":                                         faqHowToRegisterHandler,
		"Как и когда происходит оплата курса?":                                   faqWhenToPayHandler,
		"Хочу оплатить в рассрочку. Какие условия?":                              faqInstallmentHandler,
		"Я из другой страны. Могу ли я записаться на курс? Как проходит оплата?": faqForeignHandler,
		"Как понять на какой курс я могу записаться со своим уровнем?":           faqLevelHandler,
		"Возможно ли взять перерыв во время курса?":                              faqPauseHandler,
		"Выдается ли сертификат по окончании обучения?":                          faqCertificateHandler,

		// courses
		"Основы рисунка":                   courseDetailsHandler,
		"Форма и тон":                      courseDetailsHandler,
		"Свет и цвет":                      courseDetailsHandler,
		"Портрет: Скетчинг и стилизация":   courseDetailsHandler,
		"Скетчинг: тело, движение, одежда": courseDetailsHandler,
		"Динамический портрет":             courseDetailsHandler,
		"Фигура человека":                  courseDetailsHandler,
		"Мастерская с Евой":                courseDetailsHandler,
		"Дизайн существ":                   courseDetailsHandler,

		//courses details
		"Основная информация":    courseMainInfoHandler,
		"Доступные тарифы":       courseTariffsHandler,
		"Для кого подходит курс": courseForWhomHandler,
		"Программа курса":        courseScheduleHandler,
		"О чем курс":             courseAboutHandler,
		"Что понадобится":        courseToolsHandler,
		"Где посмотреть работы куратора и студентов?": WhereToFindWorksHandler,

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
		"Основы рисунка",
		"Форма и тон",
		"Свет и цвет",
		"Портрет: Скетчинг и стилизация",
		"Скетчинг: тело, движение, одежда",
		"Динамический портрет",
		"Фигура человека",
		"Мастерская с Евой",
		"Дизайн существ",
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
