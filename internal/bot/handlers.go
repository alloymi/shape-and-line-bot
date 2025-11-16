package bot

import (
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

// ===== WIP =====

func courseWIPHandler(b *Bot, m *tgbotapi.Message) {
	b.sendText(m.Chat.ID, "[этот раздел находится в разработке]")
}

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

	b.api.Send(tgbotapi.NewMessage(chatID,
		"Пожалуйста введите ваше ФИО через пробел:"))
}

func waitlistFullNameHandler(b *Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	fullname := msg.Text

	if len(fullname) < 5 || len(strings.Split(fullname, " ")) < 2 {
		b.api.Send(tgbotapi.NewMessage(chatID, "Пожалуйста, укажите ФИО полностью."))
		return
	}

	userTempFullname[chatID] = fullname
	SetState(chatID, StateWaitlistAskEmail)

	b.api.Send(tgbotapi.NewMessage(chatID,
		"Теперь введите вашу почту:"))
}

func waitlistEmailHandler(b *Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	email := msg.Text

	if !strings.Contains(email, "@") {
		b.api.Send(tgbotapi.NewMessage(chatID, "Некорректный формат почты. Попробуйте ещё раз."))
		return
	}

	course := userTempCourse[chatID]
	fullname := userTempFullname[chatID]

	// PostgreSQL
	if err := db.SaveWaitlist(b.db, chatID, fullname, email, course); err != nil {
		log.Println("DB error:", err)
		b.api.Send(tgbotapi.NewMessage(chatID, "DB save error\nПожалуйста, свяжитесь напрямую с менеджером!"))
		return
	}

	// Google Sheets
	if err := services.SaveToGoogleSheet(fullname, email, course); err != nil {
		log.Println("Sheets error:", err)
		b.api.Send(tgbotapi.NewMessage(chatID, "Sheets save error\nПожалуйста, свяжитесь напрямую с менеджером!"))
		return
	}

	ResetState(chatID)

	b.api.Send(tgbotapi.NewMessage(chatID,
		"Вы успешно записаны в лист ожидания!\nЛист ожидания не предусматривает оплаты, мы лишь уведомим вас о начале набора до официального поста в группе!\nХотим предупредить, что запись в лист ожидания не гарантирует запись на курс."))
}
