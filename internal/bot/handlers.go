package bot

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *Bot) sendText(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.api.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func menuHandler(menuName string, state BotState) Handler {
	return func(b *Bot, m *tgbotapi.Message) {
		SetState(m.Chat.ID, state)

		msg := tgbotapi.NewMessage(m.Chat.ID, "Выберите пункт:")

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

	msg := tgbotapi.NewMessage(m.Chat.ID, "Главное меню:")
	if kb, ok := Menus["main"]; ok {
		msg.ReplyMarkup = kb
	}
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Failed to send main menu: %v", err)
	}
}

// ===== FAQ ANSWERS =====
func faqHowHandler(b *Bot, m *tgbotapi.Message) {
	text := "Курс состоит из предзаписанных лекций, домашних заданий и еженедельных созвонов."
	b.sendText(m.Chat.ID, text)
}

func faqFormatHandler(b *Bot, m *tgbotapi.Message) {
	text := "Формат: видеоуроки + задания + групповые созвоны с куратором."
	b.sendText(m.Chat.ID, text)
}

func faqInstallmentHandler(b *Bot, m *tgbotapi.Message) {
	text := "Мы предлагаем рассрочку для держателей карт российских банков на 4 и 6 месяцев. Рассрочка без процентов — подробнее у менеджера."
	b.sendText(m.Chat.ID, text)
}

// ===== Courses =====

func courseWIPHandler(b *Bot, m *tgbotapi.Message) {
	b.sendText(m.Chat.ID, "[этот раздел находится в разработке]")
}

//func saveUser(b *Bot, chatID int64, username string) {
//	_, err := b.db.Exec(
//		"INSERT INTO users(chat_id, username) VALUES($1, $2) ON CONFLICT DO NOTHING",
//		chatID,
//		username,
//	)
//	if err != nil {
//		b.api.Send(tgbotapi.NewMessage(chatID, "Ошибка записи в БД"))
//		return
//	}
//}

// ===== Waiting list =====

func startWaitlistHandler(b *Bot, msg *tgbotapi.Message) {
	SetState(msg.Chat.ID, StateWaitlistChooseCourse)

	b.api.Send(tgbotapi.NewMessage(msg.Chat.ID,
		"Выберите курс, на который хотите записаться:").SetReplyMarkup(WaitlistCoursesMenu()))
}

func waitlistChooseCourseHandler(b *Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	SetState(chatID, StateWaitlistAskFullName)

	cleanName := strings.TrimPrefix(msg.Text, "WL: ")
	userTempCourse[chatID] = cleanName

	b.api.Send(tgbotapi.NewMessage(chatID,
		"Отлично! Теперь введите ваше ФИО полностью:\n\nПример: Иванова Мария Андреевна"))
}

func waitlistFullNameHandler(b *Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	full := msg.Text

	if len(full) < 5 || len(strings.Split(full, " ")) < 2 {
		b.api.Send(tgbotapi.NewMessage(chatID, "Пожалуйста, укажите ФИО полностью."))
		return
	}

	userTempFullname[chatID] = full
	SetState(chatID, StateWaitlistAskEmail)

	b.api.Send(tgbotapi.NewMessage(chatID,
		"Хорошо! Теперь введите вашу почту:\n\nПример: name@gmail.com"))
}

func waitlistEmailHandler(b *Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	email := msg.Text

	if !strings.Contains(email, "@") {
		b.api.Send(tgbotapi.NewMessage(chatID, "Почта выглядит некорректно. Попробуйте ещё раз."))
		return
	}

	course := userTempCourse[chatID]
	fullname := userTempFullname[chatID]

	// 1) PostgreSQL
	err := saveWaitlistToDB(b.db, chatID, fullname, email, course)
	if err != nil {
		b.api.Send(tgbotapi.NewMessage(chatID, "Ошибка при сохранении в БД 😢"))
		return
	}

	// 2) Google Sheets
	err = SaveToGoogleSheet(fullname, email, course)
	if err != nil {
		b.api.Send(tgbotapi.NewMessage(chatID, "Ошибка сохранения в Google Sheets"))
		return
	}

	ResetState(chatID)

	b.api.Send(tgbotapi.NewMessage(chatID,
		"Вы успешно записаны в лист ожидания! 😊\n\nМы свяжемся с вами при открытии набора."))
}
