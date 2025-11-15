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
	b.sendText(m.Chat.ID, "Курс состоит из предзаписанных лекций, домашних заданий и еженедельных созвонов.")
}

func faqFormatHandler(b *Bot, m *tgbotapi.Message) {
	b.sendText(m.Chat.ID, "Формат: видеоуроки + задания + групповые созвоны с куратором.")
}

func faqInstallmentHandler(b *Bot, m *tgbotapi.Message) {
	b.sendText(m.Chat.ID, "Мы предлагаем рассрочку на 4 и 6 месяцев. Рассрочка без процентов — подробности у менеджера.")
}

// ===== WIP =====

func courseWIPHandler(b *Bot, m *tgbotapi.Message) {
	b.sendText(m.Chat.ID, "[этот раздел находится в разработке]")
}

// ===== WAITLIST =====

func startWaitlistHandler(b *Bot, msg *tgbotapi.Message) {
	SetState(msg.Chat.ID, StateWaitlistChooseCourse)

	resp := tgbotapi.NewMessage(msg.Chat.ID, "Выберите курс, на который хотите записаться:")
	resp.ReplyMarkup = WaitlistCoursesMenu()

	b.api.Send(resp)
}

func waitlistChooseCourseHandler(b *Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	SetState(chatID, StateWaitlistAskFullName)

	// убираем префикс WL: (если есть)
	cleanName := strings.TrimPrefix(msg.Text, "WL: ")
	userTempCourse[chatID] = cleanName

	b.api.Send(tgbotapi.NewMessage(chatID,
		"Отлично! Теперь введите ваше ФИО полностью:\n\nПример: "))
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
	if err := db.SaveWaitlist(b.db, chatID, fullname, email, course); err != nil {
		log.Println("DB error:", err)
		b.api.Send(tgbotapi.NewMessage(chatID, "Ошибка при сохранении в БД"))
		return
	}

	// 2) Google Sheets
	if err := services.SaveToGoogleSheet(fullname, email, course); err != nil {
		log.Println("Sheets error:", err)
		b.api.Send(tgbotapi.NewMessage(chatID, "Ошибка сохранения в Google Sheets"))
		return
	}

	ResetState(chatID)

	b.api.Send(tgbotapi.NewMessage(chatID,
		"Вы успешно записаны в лист ожидания!"))
}
