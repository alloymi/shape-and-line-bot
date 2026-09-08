package bot

import (
	"SnLbot/internal/config"
	"SnLbot/internal/pkg/utils"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

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
		if update.CallbackQuery != nil {
			bot.processCallback(update.CallbackQuery)
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
			if update.CallbackQuery != nil {
				bot.processCallback(update.CallbackQuery)
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

//func (bot *Bot) processCallback(callback *tgbotapi.CallbackQuery) {
//	if callback == nil {
//		return
//	}
//
//	log.Printf("CALLBACK DATA: %q", callback.Data)
//
//	if callback.Message == nil {
//		return
//	}
//
//	_, _ = bot.api.Request(tgbotapi.NewCallback(callback.ID, ""))
//
//	chatID := callback.Message.Chat.ID
//	data := callback.Data
//
//	msg := *callback.Message
//	msg.Text = data
//
//	switch {
//
//	// ===== MAIN MENU =====
//
//	case data == "main":
//		startHandler(bot, &msg)
//
//	// ===== FAQ =====
//
//	case data == "faq":
//		faqHandler(bot, &msg)
//
//	case data == "faq_about":
//		msg.Text = "О школе"
//		faqCategoryHandler(bot, &msg)
//
//	case data == "faq_payment":
//		msg.Text = "Вопросы об оплате"
//		faqCategoryHandler(bot, &msg)
//
//	case data == "faq_study":
//		msg.Text = "Вопросы об обучении"
//		faqCategoryHandler(bot, &msg)
//
//	// FAQ — answers
//
//	case data == "faq_about_info":
//		faqAboutHandler(bot, &msg)
//
//		resp := tgbotapi.NewMessage(chatID, "Вопросы о школе:")
//		resp.ReplyMarkup = faqAboutSchoolInlineMenu()
//		bot.api.Send(resp)
//
//	case data == "faq_level":
//		faqLevelHandler(bot, &msg)
//
//		resp := tgbotapi.NewMessage(chatID, "Вопросы о школе:")
//		resp.ReplyMarkup = faqAboutSchoolInlineMenu()
//		bot.api.Send(resp)
//
//	case data == "faq_register":
//		faqHowToRegisterHandler(bot, &msg)
//
//		resp := tgbotapi.NewMessage(chatID, "Вопросы о школе:")
//		resp.ReplyMarkup = faqAboutSchoolInlineMenu()
//		bot.api.Send(resp)
//
//	case data == "faq_pay_when":
//		faqWhenToPayHandler(bot, &msg)
//
//		resp := tgbotapi.NewMessage(chatID, "Вопросы об оплате:")
//		resp.ReplyMarkup = faqPaymentInlineMenu()
//		bot.api.Send(resp)
//
//	case data == "faq_installment":
//		faqInstallmentHandler(bot, &msg)
//
//		resp := tgbotapi.NewMessage(chatID, "Вопросы об оплате:")
//		resp.ReplyMarkup = faqPaymentInlineMenu()
//		bot.api.Send(resp)
//
//	case data == "faq_foreign":
//		faqForeignHandler(bot, &msg)
//
//		resp := tgbotapi.NewMessage(chatID, "Вопросы об оплате:")
//		resp.ReplyMarkup = faqPaymentInlineMenu()
//		bot.api.Send(resp)
//
//	case data == "faq_format":
//		faqFormatHandler(bot, &msg)
//
//		resp := tgbotapi.NewMessage(chatID, "Вопросы об обучении:")
//		resp.ReplyMarkup = faqStudyInlineMenu()
//		bot.api.Send(resp)
//
//	case data == "faq_pause":
//		faqPauseHandler(bot, &msg)
//
//		resp := tgbotapi.NewMessage(chatID, "Вопросы об обучении:")
//		resp.ReplyMarkup = faqStudyInlineMenu()
//		bot.api.Send(resp)
//
//	case data == "faq_certificate":
//		faqCertificateHandler(bot, &msg)
//
//		resp := tgbotapi.NewMessage(chatID, "Вопросы об обучении:")
//		resp.ReplyMarkup = faqStudyInlineMenu()
//		bot.api.Send(resp)
//
//	// ===== COURSES =====
//
//	case data == "courses":
//		SetState(chatID, StateCourses)
//
//		markup := coursesInlineMenu()
//
//		edit := tgbotapi.NewEditMessageText(
//			chatID,
//			callback.Message.MessageID,
//			"Выберите курс:",
//		)
//		edit.ReplyMarkup = &markup
//
//		bot.api.Send(edit)
//
//	case strings.HasPrefix(data, "course:"):
//		course := strings.TrimPrefix(data, "course:")
//
//		msg.Text = course
//		courseDetailsHandler(bot, &msg)
//
//	// ===== COURSE INFO =====
//
//	case data == "courseinfo:main":
//		courseMainInfoHandler(bot, &msg)
//
//	case data == "courseinfo:tariffs":
//		courseTariffsHandler(bot, &msg)
//
//	case data == "courseinfo:schedule":
//		courseScheduleHandler(bot, &msg)
//
//	case data == "courseinfo:about":
//		courseAboutHandler(bot, &msg)
//
//	case data == "courseinfo:tools":
//		courseToolsHandler(bot, &msg)
//
//	case data == "courseinfo:forwhom":
//		courseForWhomHandler(bot, &msg)
//
//	case data == "courseinfo:works":
//		WhereToFindWorksHandler(bot, &msg)
//
//	// ===== WAITLIST =====
//
//	case data == "waitlist":
//		SetState(chatID, StateWaitlistChooseCourse)
//
//		markup := waitlistInlineMenu()
//
//		edit := tgbotapi.NewEditMessageText(
//			chatID,
//			callback.Message.MessageID,
//			"Выберите курс, в лист ожидания которого хотите записаться:",
//		)
//		edit.ReplyMarkup = &markup
//
//		bot.api.Send(edit)
//
//	case strings.HasPrefix(data, "waitlist:"):
//		course := strings.TrimPrefix(data, "waitlist:")
//
//		userTempCourse[chatID] = course
//
//		msg.Text = course
//		waitlistChooseCourseHandler(bot, &msg)
//	}
//}

func (bot *Bot) processCallback(callback *tgbotapi.CallbackQuery) {
	if callback == nil || callback.Message == nil {
		return
	}

	_, _ = bot.api.Request(tgbotapi.NewCallback(callback.ID, ""))

	chatID := callback.Message.Chat.ID
	data := callback.Data
	messageID := callback.Message.MessageID

	msg := *callback.Message
	msg.Text = data

	switch {

	// ===== MAIN MENU =====

	case data == "main":
		SetState(chatID, StateDefault)

		edit := tgbotapi.NewEditMessageText(
			chatID,
			messageID,
			"Здравствуйте! Это Shape and line — современная художественная онлайн-школа в Санкт-Петербурге.\n\n"+
				"В обучении мы соединили традиции классического рисунка с прогрессивными зарубежными методиками и цифровыми технологиями, чтобы вы могли учиться у лучших преподавателей, где бы вы ни находились.\n"+
				"Этот бот поможет вам сориентироваться в наших курсах, ответит на любые вопросы, а также запишет вас в лист ожидания!\n\n"+
				"Если у вас остались вопросы, вы хотите записаться на курс или нужна любая другая помощь, вы можете обратиться к нашему менеджеру @shapeandlinemanager",
		)

		markup := mainInlineMenu()
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	// ===== FAQ =====

	case data == "faq":
		SetState(chatID, StateFAQ)

		edit := tgbotapi.NewEditMessageText(
			chatID,
			messageID,
			"Выберите категорию вопросов:",
		)

		markup := faqCategoriesInlineMenu()
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	case data == "faq_about":
		SetState(chatID, StateFAQAbout)

		edit := tgbotapi.NewEditMessageText(
			chatID,
			messageID,
			"Вопросы о школе:",
		)

		markup := faqAboutSchoolInlineMenu()
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	case data == "faq_payment":
		SetState(chatID, StateFAQPayment)

		edit := tgbotapi.NewEditMessageText(
			chatID,
			messageID,
			"Вопросы об оплате:",
		)

		markup := faqPaymentInlineMenu()
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	case data == "faq_study":
		SetState(chatID, StateFAQStudy)

		edit := tgbotapi.NewEditMessageText(
			chatID,
			messageID,
			"Вопросы об обучении:",
		)

		markup := faqStudyInlineMenu()
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	// ===== FAQ ANSWERS =====

	case data == "faq_about_info":
		edit := tgbotapi.NewEditMessageText(
			chatID,
			messageID,
			"Shape and Line — онлайн-школа рисования из Санкт-Петербурга.\n\n"+
				"Мы объединяем классическую академическую базу, современные цифровые технологии "+
				"и опыт кураторов, которые работают в игровой индустрии, анимации и комиксах.\n\n"+
				"На курсах мы разбираем, как устроены форма, свет, цвет, перспектива, анатомия, "+
				"композиция и многое другое. Это та основа, которая остаётся с вами независимо от того, "+
				"в каком направлении вы решите развиваться дальше!",
		)

		markup := faqAboutSchoolInlineMenu()
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	case data == "faq_level":
		edit := tgbotapi.NewEditMessageText(
			chatID,
			messageID,
			"Мы будем рады помочь вам с выбором!\n\n"+
				"Чтобы оценить ваш уровень и подобрать подходящий курс, "+
				"пожалуйста, расскажите о ваших актуальных целях. Также вы можете прислать небольшое портфолио из 4–5 работ ссылкой на диск, "+
				"вашу группу или файлами — это поможет нашим кураторам подобрать для вас курс, который будет для вас сейчас наиболее полезным.\n\n"+
				"Всю подготовленную информацию и портфолио вы можете отправить нашему менеджеру @shapeandlinemanager — и мы поможем подобрать для вас курс!",
		)

		markup := faqAboutSchoolInlineMenu()
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	case data == "faq_register":
		edit := tgbotapi.NewEditMessageText(
			chatID,
			messageID,
			"Для записи на курс можете обратиться к нашему менеджеру: @shapeandlinemanager\n\n"+
				"Для записи вас на курс потребуется подготовить портфолио ваших актуальных работ. Это актуально для всех курсов, кроме «Основ рисунка», "+
				"так как прием в группу осуществляется только после одобрения ваших работ куратором.\n"+
				"Это может быть ссылка на артстейшн, сообщество в соцсетях или на папку с работами на гугл-диске, отражающих ваш уровень!",
		)

		markup := faqAboutSchoolInlineMenu()
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	case data == "faq_pay_when":
		edit := tgbotapi.NewEditMessageText(
			chatID,
			messageID,
			"Оплату курса необходимо произвести до начала курса в любое удобное для вас время."+
				"Но мы рекомендуем не затягивать, так как количество мест в группе ограничено, а бронь места возможна только при оплате!\n\n"+
				"Ссылку на оплату вы получите только после подписания договора. Наш куратор подготовит для вас персональную ссылку для оплаты.",
		)

		markup := faqPaymentInlineMenu()
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	case data == "faq_installment":
		edit := tgbotapi.NewEditMessageText(
			chatID,
			messageID,
			"Мы предлагаем рассрочку для держателей карт российских банков на 4 и 6 месяцев. Рассрочка без процентов и предоставляется от Т-банка.\n"+
				"*Банк вправе установить комиссию (проценты) для Клиентов за предоставление рассрочки на приобретение Товара или иную ставку, в связи с чем у"+
				"Клиента может возникнуть переплата за Товар. Банк самостоятельно определяет размер комиссии (процентов) и повышенной ставки и иные условия их"+
				"расчета и начисления по своему усмотрению.\n\nЕщё оплату можно внести долями. Подробнее о сервисе Долями по ссылке: https://dolyame.ru/help/customer/about/",
		)

		markup := faqPaymentInlineMenu()
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	case data == "faq_foreign":
		edit := tgbotapi.NewEditMessageText(
			chatID,
			messageID,
			"Мы принимаем оплату из других стран переводом куратору через сервис PayPal!\n"+
				"Если данный способ вам не подходит, вы можете уточнить варианты оплаты у менеджера @shapeandlinemanager",
		)

		markup := faqPaymentInlineMenu()
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	case data == "faq_format":
		edit := tgbotapi.NewEditMessageText(
			chatID,
			messageID,
			"Курсы представлены в формате предзаписанных лекций, в конце которых содержится домашнее задание. "+
				"Просматриваете и выполняете задания вы самостоятельно.  Лекции предоставляются в формате файлов для скачивания, которые "+
				"доступны для просмотра через Инфопротектор. Доступ к лекционным материалам предоставляется студентам бессрочно.\n\n"+
				"Раз в неделю в определённое время проходит групповой созвон, где вы получаете фидбек на домашнее задание от куратора. "+
				"Созвоны в основном проходят в 19:00 по МСК, так же у студентов есть доступ к записям фидбеков.",
		)

		markup := faqStudyInlineMenu()
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	case data == "faq_pause":
		edit := tgbotapi.NewEditMessageText(
			chatID,
			messageID,
			"В случае непредвиденных обстоятельств или отпуска вы можете взять перерыв на некоторое время, "+
				"но разборы домашних заданий будут идти в обычном режиме.\nВы можете догнать группу, но куратор не сможет разобрать ваши домашние задания "+
				"с пропущенных недель в полном объёме!",
		)

		markup := faqStudyInlineMenu()
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	case data == "faq_certificate":
		edit := tgbotapi.NewEditMessageText(
			chatID,
			messageID,
			"Мы предоставляем электронный сертификат об успешном завершении курса по вашему запросу!\n\n"+
				"Но уточним, что это не диплом о профессиональной переподготовке и не официальный сертификат о повышении квалификации.",
		)

		markup := faqStudyInlineMenu()
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	// ===== COURSES =====

	case data == "courses":
		SetState(chatID, StateCourses)

		edit := tgbotapi.NewEditMessageText(
			chatID,
			messageID,
			"Выберите курс:",
		)

		markup := coursesInlineMenu()
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	case strings.HasPrefix(data, "course:"):
		course := strings.TrimPrefix(data, "course:")

		userTempCourse[chatID] = course
		SetState(chatID, StateCourseMenu)

		info := CoursesInfo[course]

		// Если текущее сообщение — текстовое меню, превращаем его в сообщение с курсом.
		if callback.Message.Photo == nil {
			photo := tgbotapi.NewPhoto(
				chatID,
				tgbotapi.FileURL(info.ImageURL),
			)

			photo.Caption = fmt.Sprintf(
				"Что вы хотите узнать о курсе «%s»?",
				course,
			)

			photo.ReplyMarkup = courseInlineMenu(course)

			// Удаляем старое меню, чтобы сообщения не копились.
			deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
			bot.api.Send(deleteMsg)

			bot.api.Send(photo)
		} else {
			edit := tgbotapi.NewEditMessageCaption(
				chatID,
				messageID,
				fmt.Sprintf(
					"Что вы хотите узнать о курсе «%s»?",
					course,
				),
			)

			markup := courseInlineMenu(course)
			edit.ReplyMarkup = &markup
			bot.api.Send(edit)
		}

	// ===== COURSE INFO =====

	case data == "courses":
		SetState(chatID, StateCourses)

		// Удаляем старое сообщение с главным меню
		deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
		bot.api.Send(deleteMsg)

		// Создаём меню курсов
		newMsg := tgbotapi.NewMessage(chatID, "Выберите курс:")
		markup := coursesInlineMenu()
		newMsg.ReplyMarkup = markup

		bot.api.Send(newMsg)
	//case data == "courseinfo:main":
	//	course := userTempCourse[chatID]
	//	info := CoursesInfo[course]
	//
	//	edit := tgbotapi.NewEditMessageCaption(
	//		chatID,
	//		messageID,
	//		info.MainInfo,
	//	)
	//
	//	markup := courseInlineMenu(course)
	//	edit.ReplyMarkup = &markup
	//	bot.api.Send(edit)

	case data == "courseinfo:tariffs":
		course := userTempCourse[chatID]
		info := CoursesInfo[course]

		edit := tgbotapi.NewEditMessageCaption(
			chatID,
			messageID,
			info.Tariffs,
		)

		markup := courseInlineMenu(course)
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	case data == "courseinfo:schedule":
		course := userTempCourse[chatID]
		info := CoursesInfo[course]

		edit := tgbotapi.NewEditMessageCaption(
			chatID,
			messageID,
			info.Schedule,
		)

		markup := courseInlineMenu(course)
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	case data == "courseinfo:about":
		course := userTempCourse[chatID]
		info := CoursesInfo[course]

		edit := tgbotapi.NewEditMessageCaption(
			chatID,
			messageID,
			info.About,
		)

		markup := courseInlineMenu(course)
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	case data == "courseinfo:tools":
		course := userTempCourse[chatID]
		info := CoursesInfo[course]

		edit := tgbotapi.NewEditMessageCaption(
			chatID,
			messageID,
			info.Tools,
		)

		markup := courseInlineMenu(course)
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	case data == "courseinfo:forwhom":
		course := userTempCourse[chatID]
		info := CoursesInfo[course]

		edit := tgbotapi.NewEditMessageCaption(
			chatID,
			messageID,
			info.ForWhom,
		)

		markup := courseInlineMenu(course)
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	case data == "courseinfo:works":
		course := userTempCourse[chatID]
		info := CoursesInfo[course]

		edit := tgbotapi.NewEditMessageCaption(
			chatID,
			messageID,
			info.WhereToFindWorks,
		)

		markup := courseInlineMenu(course)
		edit.ReplyMarkup = &markup
		bot.api.Send(edit)

	// ===== WAITLIST =====

	case data == "waitlist":
		SetState(chatID, StateWaitlistChooseCourse)

		// Удаляем старое сообщение с главным меню
		deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
		bot.api.Send(deleteMsg)

		// Создаём меню листа ожидания
		newMsg := tgbotapi.NewMessage(
			chatID,
			"Лист ожидания не предусматривает оплаты, мы лишь уведомим вас о начале набора до официального поста в группе!\n"+
				"Хотим предупредить, что запись в лист ожидания не гарантирует запись на курс.\n\n"+
				"Выберите курс, на который хотите записаться в лист ожидания:",
		)

		markup := waitlistInlineMenu()
		newMsg.ReplyMarkup = markup

		bot.api.Send(newMsg)
	//case data == "waitlist":
	//	SetState(chatID, StateWaitlistChooseCourse)
	//
	//	edit := tgbotapi.NewEditMessageText(
	//		chatID,
	//		messageID,
	//		"Лист ожидания не предусматривает оплаты, мы лишь уведомим вас о начале набора до официального поста в группе!\n"+
	//			"Хотим предупредить, что запись в лист ожидания не гарантирует запись на курс.\n\n"+
	//			"Выберите курс, на который хотите записаться в лист ожидания:",
	//	)
	//
	//	markup := waitlistInlineMenu()
	//	edit.ReplyMarkup = &markup
	//	bot.api.Send(edit)

	case strings.HasPrefix(data, "waitlist:"):
		course := strings.TrimPrefix(data, "waitlist:")

		userTempCourse[chatID] = course
		SetState(chatID, StateWaitlistAskFullName)

		edit := tgbotapi.NewEditMessageText(
			chatID,
			messageID,
			"Пожалуйста введите ваше ФИО через пробел:",
		)

		markup := WaitlistProgressMenu()

		edit.ReplyMarkup = nil

		bot.api.Send(edit)

		_ = markup

	}
}
