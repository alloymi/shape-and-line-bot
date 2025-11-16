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

func startHandler(b *Bot, m *tgbotapi.Message) {
	SetState(m.Chat.ID, StateDefault)

	msg := tgbotapi.NewMessage(m.Chat.ID, "Здравствуйте! Это Shape and line – современная художественная школа в Санкт-Петербурге с онлайн обучением. \n\nДля максимального результата совмещаем цифровые технологии, прогрессивные зарубежные методики и опыт русского классического образования.\n\nБот расскажет вам про наши курсы, проконсультирует насчет актуальных вопросов, а так же поможет с бронированием и записью в список ожидания!")
	if kb, ok := Menus["main"]; ok {
		msg.ReplyMarkup = kb
	}

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Failed to send main menu: %v", err)
	}
}

// ===== FAQ ANSWERS =====

func faqHowHandler(bot *Bot, m *tgbotapi.Message) {
	bot.sendText(m.Chat.ID, "Курс состоит из предзаписанных лекций, к которым мы выдаём вам доступ и которые вы отсмаьриваете самостоятельно; домашнего задания; и групповых созвонов с куратором раз в неделю, где он даёт детальный фидбек на вашу работу!")
}

func faqFormatHandler(bot *Bot, m *tgbotapi.Message) {
	bot.sendText(m.Chat.ID, "Курсы представлены в формате предзаписанных лекций, в конце которых содержится домашнее задание. Просматриваете и выполняете задания вы самостоятельно.\n\nЛекции предоставляются в формате файлов для скачивания, которые доступны для просмотра через Инфопротектор. Доступ к лекционным материалам предоставляется студентам бессрочно.\n\nРаз в неделю в определённое время проходит групповой созвон, где вы получаете фидбек на домашнее задание от куратора. Созвоны в основном проходят в 19:00 по МСК, так же у студентов есть доступ к записям фидбеков.\"")
}

func faqInstallmentHandler(bot *Bot, m *tgbotapi.Message) {
	bot.sendText(m.Chat.ID, "Мы предлагаем рассрочку для держателей карт российских банков на 4 и 6 месяцев. Рассрочка без процентов и предоставляется от Т-банка!\n")
}

func faqForeignHandler(bot *Bot, m *tgbotapi.Message) {
	bot.sendText(m.Chat.ID, "Мы принимаем оплату из других стран переводом куратору через сервис PayPal!\nЕсли такой способ оплаты не подходит, можете уточнить какие есть еще варианты у администратора")
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

	msg := tgbotapi.NewMessage(chatID,
		fmt.Sprintf("Что вы хотите узнать о курсе «%s»?", course))
	msg.ReplyMarkup = CourseMenu(course)

	b.api.Send(msg)
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

	resp := tgbotapi.NewMessage(msg.Chat.ID, "Выберите курс, на который хотите записаться в лист ожидания:")
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

	if fullname == "Отменить процесс записи" || fullname == "Назад в главное меню" {
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

	resp := tgbotapi.NewMessage(chatID, "Теперь введите вашу почту:")
	resp.ReplyMarkup = WaitlistProgressMenu()
	b.api.Send(resp)
}

func waitlistEmailHandler(bot *Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	email := msg.Text

	if email == "Отменить процесс записи" || email == "Назад в главное меню" {
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
		"Вы успешно записаны в лист ожидания!\n\nЛист ожидания не предусматривает оплаты, мы лишь уведомим вас о начале набора до официального поста в группе!\nХотим предупредить, что запись в лист ожидания не гарантирует запись на курс."))

	ResetState(chatID)

	mainMenuMsg := tgbotapi.NewMessage(chatID, "Возвращение в главное меню:")
	if kb, ok := Menus["main"]; ok {
		mainMenuMsg.ReplyMarkup = kb
	}
	bot.api.Send(mainMenuMsg)
}
