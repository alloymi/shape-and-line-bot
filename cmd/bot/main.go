package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Panic("Bot token not found in .env")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleMessage(bot, update.Message)
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID

	switch message.Text {
	case "/start":
		sendMainMenu(bot, chatID)
	case "/help":
		sendMainMenu(bot, chatID)

	//main menu
	case "частые вопросы":
		sendFAQMenu(bot, chatID)
	case "вопрос о курсе":
		sendCoursesMenu(bot, chatID)
	case ":D":
		sendText(bot, chatID, "D:")

	// faq menu
	case "Как проходит обучение?":
		sendText(bot, chatID, "Курс состоит из предзаписанных лекций, к которым мы выдаём вам доступ и которые вы отсматриваете самостоятельно; домашнего задания; и групповых созвонов с куратором раз в неделю, где он даёт детальный фидбек на вашу работу!")
	case "Формат обучения":
		sendText(bot, chatID, "\"Курсы представлены в формате предзаписанных лекций, в конце которых содержится домашнее задание. Просматриваете и выполняете задания вы самостоятельно.\n\nЛекции предоставляются в формате файлов для скачивания, которые доступны для просмотра через Инфопротектор. Доступ к лекционным материалам предоставляется студентам бессрочно.\n\nРаз в неделю в определённое время проходит групповой созвон, где вы получаете фидбек на домашнее задание от куратора. Созвоны в основном проходят в 19:00 по МСК, так же у студентов есть доступ к записям фидбеков.")
	case "Хочу оплатить в рассрочку. Какие условия?":
		sendText(bot, chatID, "Мы предлагаем рассрочку для держателей карт российских банков на 4 и 6 месяцев. Рассрочка без процентов и предоставляется от Т-банка!")
	case "okak":
		sendText(bot, chatID, "avottak")
	case "назад":
		sendMainMenu(bot, chatID)

	// courses menu
	case "назад к выбору курса":
		sendCoursesMenu(bot, chatID)
	case "фигура человека":
		sendText(bot, chatID, "тут пока ничего нет")
	case "форма и тон":
		sendText(bot, chatID, "тут пока ничего нет")
	case "дизайн существ":
		sendText(bot, chatID, "тут пока ничего нет")
	case "портрет: скетчинг и стилизация":
		sendText(bot, chatID, "тут пока ничего нет")
	case "свет и цвет":
		sendText(bot, chatID, "тут пока ничего нет")
	case "диманический портрет":
		sendText(bot, chatID, "тут пока ничего нет")
	case "основы рисунка":
		sendText(bot, chatID, "тут пока ничего нет")
	case "мастерская с Евой":
		sendText(bot, chatID, "тут пока ничего нет")
	case "анатомия человека":
		sendText(bot, chatID, "тут пока ничего нет")

	default:
		sendMainMenu(bot, chatID)
	}
}

func sendMainMenu(bot *tgbotapi.BotAPI, chatID int64) {
	text := "[тут какое-то приветсвие]"

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("частые вопросы"),
			tgbotapi.NewKeyboardButton("вопрос о курсе"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(":D"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func sendFAQMenu(bot *tgbotapi.BotAPI, chatID int64) {
	text := "выберите вопрос:"

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Как проходит обучение?"),
			tgbotapi.NewKeyboardButton("Формат обучения"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Хочу оплатить в рассрочку. Какие условия?"),
			tgbotapi.NewKeyboardButton("okak"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("назад"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func sendText(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}

func sendCoursesMenu(bot *tgbotapi.BotAPI, chatID int64) {
	text := "выберите курс по которому хотите задать вопрос:"

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("фигура человека"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("форма и тон"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("дизайн существ"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("портрет: скечтинг и стилизация"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("свет и цвет"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("диманический портрет"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("основы рисунка"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("мастерская с Евой"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("анатомия человека"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("назад"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}
