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
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Записаться в лист ожидания")),
	)
}

func faqMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Как проходит обучение?"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Формат обучения?"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Хочу оплатить в рассрочку. Какие условия?"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Я из другой страны. Могу ли я записаться на курс? Как проходит оплата?"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад в главное меню"),
		),
	)
}

// courses

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
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Назад в главное меню")))
}

//func CourseDetailsMenu() tgbotapi.ReplyKeyboardMarkup {
//	return tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Длительность курса")),
//		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Ближайший старт")),
//		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Куратор курса")),
//		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Назад к списку курсов")),
//	)
//}

func CourseMenu(courseName string) tgbotapi.ReplyKeyboardMarkup {
	info := CoursesInfo[courseName]

	buttons := [][]tgbotapi.KeyboardButton{}

	if info.Duration != "" {
		buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Длительность курса"),
		))
	}

	if info.StartDate != "" {
		buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Ближайший старт"),
		))
	}

	if info.Curator != "" {
		buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Куратор курса"),
		))
	}

	if info.Schedule != "" {
		buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Программа курса"),
		))
	}

	if info.About != "" {
		buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("О чем курс"),
		))
	}

	if info.Tools != "" {
		buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Что понадобится"),
		))
	}

	buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Назад к списку курсов"),
	))

	return tgbotapi.NewReplyKeyboard(buttons...)
}

// waitlist

func WaitlistCoursesMenu() tgbotapi.ReplyKeyboardMarkup {
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
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Назад в главное меню")))
}

func WaitlistProgressMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отменить процесс записи"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад в главное меню"),
		),
	)
}
