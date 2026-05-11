package bot

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"SnLbot/internal/db"
	"SnLbot/internal/services"
)

func (bot *Bot) sendText(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, _ = bot.api.Send(msg)
}

func menuHandler(menuName string, state BotState) Handler {
	return func(b *Bot, m *tgbotapi.Message) {
		SetState(m.Chat.ID, state)

		msg := tgbotapi.NewMessage(m.Chat.ID, "Выберите интересующий вас пункт:")

		if kb, ok := Menus[menuName]; ok {
			msg.ReplyMarkup = kb
		}

		if _, err := b.api.Send(msg); err != nil {
			log.Printf("Failed to send menu %s: %v", menuName, err)
		}
	}
}

func startHandler(bot *Bot, m *tgbotapi.Message) {
	SetState(m.Chat.ID, StateDefault)

	photoUrl := "https://i.postimg.cc/g00BzdVp/logobot.png"

	msg := tgbotapi.NewMessage(m.Chat.ID, "Здравствуйте! Это Shape and line — современная художественная онлайн-школа в Санкт-Петербурге.\n\n"+
		"В обучении мы соединили традиции классического рисунка с прогрессивными зарубежными методиками и цифровыми технологиями, чтобы вы могли учиться у лучших преподавателей, где бы вы ни находились.\n"+
		"Этот бот поможет вам сориентироваться в наших курсах, ответит на любые вопросы, а также запишет вас в лист ожидания!\n\n"+
		"Если у вас остались вопросы, вы хотите записаться на курс или нужна любая другая помощь, вы можете обратиться к нашему менеджеру @shapeandlinemanager")

	photo := tgbotapi.NewPhoto(m.Chat.ID, tgbotapi.FileURL(photoUrl))
	_, err := bot.api.Send(photo)
	if err != nil {
		log.Printf("Failed to send welcome image: %v", err)
	}

	if kb, ok := Menus["main"]; ok {
		msg.ReplyMarkup = kb
	}

	if _, err := bot.api.Send(msg); err != nil {
		log.Printf("Failed to send main menu: %v", err)
	}
}

// ===== FAQ ANSWERS =====

func faqAboutHandler(bot *Bot, m *tgbotapi.Message) {
	bot.sendText(m.Chat.ID, "Shape and Line — онлайн-школа рисования из Санкт-Петербурга.\n\n"+
		"Мы объединяем классическую академическую базу, современные цифровые технологии "+
		"и опыт кураторов, которые работают в игровой индустрии, анимации и комиксах.\n\n"+
		"На курсах мы разбираем, как устроены форма, свет, цвет, перспектива, анатомия, "+
		"композиция и многое другое. Это та основа, которая остаётся с вами независимо от того, "+
		"в каком направлении вы решите развиваться дальше!")
}

//func faqHowHandler(bot *Bot, m *tgbotapi.Message) {
//	bot.sendText(m.Chat.ID, "Курс состоит из предзаписанных лекций, к которым мы выдаём вам доступ и которые вы отсматриваете самостоятельно; домашнего задания; и групповых созвонов с куратором раз в неделю, где он даёт детальный фидбек на вашу работу!")
//}

func faqFormatHandler(bot *Bot, m *tgbotapi.Message) {
	bot.sendText(m.Chat.ID, "Курсы представлены в формате предзаписанных лекций, в конце которых содержится домашнее задание. Просматриваете и выполняете задания вы самостоятельно.\n\nЛекции предоставляются в формате файлов для скачивания, которые доступны для просмотра через Инфопротектор. Доступ к лекционным материалам предоставляется студентам бессрочно.\n\nРаз в неделю в определённое время проходит групповой созвон, где вы получаете фидбек на домашнее задание от куратора. Созвоны в основном проходят в 19:00 по МСК, так же у студентов есть доступ к записям фидбеков.")
}

func faqHowToRegisterHandler(bot *Bot, m *tgbotapi.Message) {
	bot.sendText(m.Chat.ID, "Для записи на курс можете обратиться к нашему менеджеру: @shapeandlinemanager\n\nДля записи вас на курс потребуется подготовить портфолио ваших актуальных работ. Это актуально для всех курсов, кроме «Основ рисунка», так как прием в группу осуществляется только после одобрения ваших работ куратором.\nЭто может быть ссылка на артстейшн, сообщество в соцсетях или на папку с работами на гугл-диске, отражающих ваш уровень!")
}

func faqWhenToPayHandler(bot *Bot, m *tgbotapi.Message) {
	bot.sendText(m.Chat.ID, "Оплату курса необходимо произвести до начала курса в любое удобное для вас время. Но мы рекомендуем не затягивать, так как количество мест в группе ограничено, а бронь места возможна только при оплате!\n\nСсылку на оплату вы получите только после подписания договора. Наш куратор подготовит для вас персональную ссылку для оплаты.")
}

func faqInstallmentHandler(bot *Bot, m *tgbotapi.Message) {
	bot.sendText(m.Chat.ID, "Мы предлагаем рассрочку для держателей карт российских банков на 4 и 6 месяцев. Рассрочка без процентов и предоставляется от Т-банка. \n*Банк вправе установить комиссию (проценты) для Клиентов за предоставление рассрочки на приобретение Товара\nили иную ставку, в связи с чем у Клиента может возникнуть переплата за Товар. Банк самостоятельно определяет размер комиссии (процентов) и повышенной ставки и иные условия их расчета и начисления по своему усмотрению.\n\n\nЕщё оплату можно внести долями. Подробнее о сервисе Долями по ссылке: https://dolyame.ru/help/customer/about/")
}

func faqForeignHandler(bot *Bot, m *tgbotapi.Message) {
	bot.sendText(m.Chat.ID, "Мы принимаем оплату из других стран переводом куратору через сервис PayPal!\nЕсли данный способ вам не подходит, вы можете уточнить варианты оплаты у менеджера @shapeandlinemanager")
}

func faqLevelHandler(bot *Bot, m *tgbotapi.Message) {
	bot.sendText(m.Chat.ID, "ы будем рады помочь вам с выбором!\n\nЧтобы оценить ваш уровень и подобрать подходящий курс, расскажите о ваших актуальных целях. Также вы можете прислать небольшое портфолио из 4–5 работ ссылкой на диск, вашу группу или файлами — это поможет нашим кураторам подобрать для вас курс, который будет для вас сейчас наиболее полезным.\n\nВсю подготовленную информацию и портфолио вы можете отправить нашему менеджеру @shapeandlinemanager — и мы поможем подобрать для вас курс!")
}

func faqPauseHandler(bot *Bot, m *tgbotapi.Message) {
	bot.sendText(m.Chat.ID, "В случае непредвиденных обстоятельств или, например, отпуска на некоторое время перерыв взять можно, но разборы домашних заданий будут идти в обычном режиме. Вы сможете догнать группу, но имейте в виду, что куратор не сможет разобрать ваши домашние задания с пропущенных недель в полном объёме!")
}

func faqCertificateHandler(bot *Bot, m *tgbotapi.Message) {
	bot.sendText(m.Chat.ID, "Мы предоставляем электронный сертификат об успешном завершении курса по вашему запросу!\n\nНо уточним, что это не диплом о профессиональной переподготовке и не официальный сертификат о повышении квалификации.")
}

// ===== COURSES =====

//func courseWIPHandler(b *Bot, m *tgbotapi.Message) {
//	b.sendText(m.Chat.ID, "[этот раздел находится в разработке]")
//}

//func courseDetailsHandler(b *Bot, m *tgbotapi.Message) {
//	chatID := m.Chat.ID
//
//	userTempCourse[chatID] = m.Text
//
//	SetState(chatID, StateCourseDetails)
//
//	msg := tgbotapi.NewMessage(chatID,
//		fmt.Sprintf("Что вы хотите узнать о курсе «%s»?", m.Text))
//	msg.ReplyMarkup = CourseDetailsMenu()
//	b.api.Send(msg)
//}

func courseDetailsHandler(b *Bot, m *tgbotapi.Message) {
	chatID := m.Chat.ID
	course := m.Text

	// запоминаем выбранный курс
	userTempCourse[chatID] = course
	SetState(chatID, StateCourseMenu)

	img, _ := CoursesInfo[course]

	imgMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(img.ImageURL))
	imgMsg.Caption = fmt.Sprintf("Что вы хотите узнать о курсе «%s»?", course)
	imgMsg.ReplyMarkup = CourseMenu(course)

	b.api.Send(imgMsg)

	//msg := tgbotapi.NewMessage(chatID,
	//	fmt.Sprintf("Что вы хотите узнать о курсе «%s»?", course))
	//msg.ReplyMarkup = CourseMenu(course)
	//
	//b.api.Send(msg)
}

//func courseDurationHandler(b *Bot, m *tgbotapi.Message) {
//	course := userTempCourse[m.Chat.ID]
//
//	switch course {
//	case "Фигура человека":
//		b.sendText(m.Chat.ID, "Длительность: ")
//	case "Форма и тон":
//		b.sendText(m.Chat.ID, "Длительность: ")
//	default:
//		b.sendText(m.Chat.ID, "Информация пока недоступна. Можете обратиться к администратору, чтобы получить ответ на свой вопрос!")
//	}
//}
//
//func courseStartHandler(b *Bot, m *tgbotapi.Message) {
//	course := userTempCourse[m.Chat.ID]
//
//	switch course {
//	case "Фигура человека":
//		b.sendText(m.Chat.ID, "Ближайший старт: ")
//	case "Форма и тон":
//		b.sendText(m.Chat.ID, "Ближайший старт: ")
//	default:
//		b.sendText(m.Chat.ID, "Информация пока недоступна. Можете обратиться к администратору, чтобы получить ответ на свой вопрос!")
//	}
//}
//
//func courseTeacherHandler(b *Bot, m *tgbotapi.Message) {
//	course := userTempCourse[m.Chat.ID]
//
//	switch course {
//	case "Фигура человека":
//		b.sendText(m.Chat.ID, "Куратор: ")
//	case "Форма и тон":
//		b.sendText(m.Chat.ID, "Куратор: ")
//	default:
//		b.sendText(m.Chat.ID, "Информация пока недоступна.")
//	}
//}
//
//func courseBackHandler(b *Bot, m *tgbotapi.Message) {
//	delete(userTempCourse, m.Chat.ID)
//	SetState(m.Chat.ID, StateCourses)
//
//	msg := tgbotapi.NewMessage(m.Chat.ID, "Выберите курс:")
//	msg.ReplyMarkup = Menus["courses"]
//	b.api.Send(msg)
//}

// ===== WAITLIST =====

func startWaitlistHandler(b *Bot, msg *tgbotapi.Message) {
	SetState(msg.Chat.ID, StateWaitlistChooseCourse)

	resp := tgbotapi.NewMessage(msg.Chat.ID, "Лист ожидания не предусматривает оплаты, мы лишь уведомим вас о начале набора до официального поста в группе!\n"+
		"Хотим предупредить, что запись в лист ожидания не гарантирует запись на курс.\n\n"+
		"Выберите курс, на который хотите записаться в лист ожидания:")

	resp.ReplyMarkup = WaitlistCoursesMenu()

	b.api.Send(resp)
}

func waitlistChooseCourseHandler(b *Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	SetState(chatID, StateWaitlistAskFullName)

	//cleanName := strings.TrimPrefix(msg.Text, "WL:")
	//userTempCourse[chatID] = cleanName
	userTempCourse[chatID] = msg.Text

	resp := tgbotapi.NewMessage(chatID, "Пожалуйста введите ваше ФИО через пробел:")
	resp.ReplyMarkup = WaitlistProgressMenu()

	b.api.Send(resp)

}

func waitlistFullNameHandler(b *Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	fullname := msg.Text

	if fullname == "Отменить процесс записи" {
		SetState(chatID, StateWaitlistChooseCourse)

		resp := tgbotapi.NewMessage(chatID, "Выберите курс, на который хотите записаться в лист ожидания:")
		resp.ReplyMarkup = WaitlistCoursesMenu()
		b.api.Send(resp)
		return
	}

	if fullname == "Назад в главное меню" {
		resetToMainMenu(b, chatID)
		return
	}

	if isCourseName(fullname) {
		b.api.Send(tgbotapi.NewMessage(chatID,
			"Пожалуйста, введите ваше ФИО:"))
		return
	}

	if len(fullname) < 5 || len(strings.Split(fullname, " ")) < 2 {
		b.api.Send(tgbotapi.NewMessage(chatID, "Пожалуйста, укажите ФИО полностью."))
		return
	}

	userTempFullname[chatID] = fullname
	SetState(chatID, StateWaitlistAskEmail)

	resp := tgbotapi.NewMessage(chatID, "Хорошо! Теперь введите вашу почту:\n\nПример: name@gmail.com")
	resp.ReplyMarkup = WaitlistProgressMenu()
	b.api.Send(resp)
}

func waitlistEmailHandler(bot *Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	email := msg.Text

	if email == "Отменить процесс записи" {
		SetState(chatID, StateWaitlistChooseCourse)

		resp := tgbotapi.NewMessage(chatID, "Выберите курс, на который хотите записаться в лист ожидания:")

		resp.ReplyMarkup = WaitlistCoursesMenu()
		bot.api.Send(resp)
		return
	}

	if email == "Назад в главное меню" {
		resetToMainMenu(bot, chatID)
		return
	}

	if !strings.Contains(email, "@") {
		bot.api.Send(tgbotapi.NewMessage(chatID, "Некорректный формат почты. Попробуйте ещё раз."))
		return
	}

	course := userTempCourse[chatID]
	fullname := userTempFullname[chatID]

	// PostgreSQL
	if err := db.SaveWaitlist(bot.db, chatID, fullname, email, course); err != nil {
		log.Println("DB error:", err)
		bot.api.Send(tgbotapi.NewMessage(chatID, "DB save error\nПожалуйста, свяжитесь напрямую с менеджером!"))
		return
	}

	// Google Sheets
	if err := services.SaveToGoogleSheet(fullname, email, course); err != nil {
		log.Println("Sheets error:", err)
		bot.api.Send(tgbotapi.NewMessage(chatID, "Sheets save error\nПожалуйста, свяжитесь напрямую с менеджером!"))
		return
	}

	summary := fmt.Sprintf(
		"Ваши данные:\n\nФИО:  %s \nПочта:  %s \nКурс:  %s",
		fullname, email, course,
	)
	bot.api.Send(tgbotapi.NewMessage(chatID, summary))

	bot.api.Send(tgbotapi.NewMessage(chatID,
		"Вы успешно записаны в лист ожидания!\n"))

	ResetState(chatID)

	mainMenuMsg := tgbotapi.NewMessage(chatID, "Возвращение в главное меню:")
	if kb, ok := Menus["main"]; ok {
		mainMenuMsg.ReplyMarkup = kb
	}
	bot.api.Send(mainMenuMsg)
}
