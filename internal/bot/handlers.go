package bot

import (
	"log"

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
