package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func courseDurationHandler(bot *Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	course := userTempCourse[chatID]
	info := CoursesInfo[course]

	bot.api.Send(tgbotapi.NewMessage(chatID, info.Duration))
}

func courseStartHandler(bot *Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	course := userTempCourse[chatID]
	info := CoursesInfo[course]

	//if info.StartDate == "" {
	//	bot.api.Send(tgbotapi.NewMessage(chatID, "Информация о старте пока недоступна"))
	//	return
	//}

	bot.api.Send(tgbotapi.NewMessage(chatID, info.StartDate))
}

func courseTeacherHandler(bot *Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	course := userTempCourse[chatID]
	info := CoursesInfo[course]

	bot.api.Send(tgbotapi.NewMessage(chatID, info.Curator))
}

func courseScheduleHandler(bot *Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	course := userTempCourse[chatID]
	info := CoursesInfo[course]

	bot.api.Send(tgbotapi.NewMessage(chatID, info.Schedule))
}

func courseAboutHandler(bot *Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	course := userTempCourse[chatID]
	info := CoursesInfo[course]

	bot.api.Send(tgbotapi.NewMessage(chatID, info.About))
}

func courseToolsHandler(bot *Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	course := userTempCourse[chatID]
	info := CoursesInfo[course]

	bot.api.Send(tgbotapi.NewMessage(chatID, info.Tools))
}

func courseBackHandler(bot *Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	ResetState(chatID)

	resp := tgbotapi.NewMessage(chatID, "Выберите курс:")
	resp.ReplyMarkup = Menus["courses"]
	bot.api.Send(resp)
}
