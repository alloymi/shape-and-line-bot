package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var Menus = map[string]tgbotapi.ReplyKeyboardMarkup{
	"main":    mainMenu(),
	"faq":     faqMenu(),
	"courses": coursesMenu(),
}

func mainMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Частые вопросы"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Все курсы"),
		),
	)
}

func faqMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Как проходит обучение?"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Формат обучения"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("назад"),
		),
	)
}

func coursesMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Фигура человека")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Форма и тон")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Дизайн существ")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Портрет: Скетчинг и стилизация")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Свет и цвет")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Динамический портрет")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Основы рисунка")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Мастерская с Евой")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Анатомия человека")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("назад")))
}
