package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var Menus = map[string]tgbotapi.ReplyKeyboardMarkup{
	"main":          mainMenu(),
	"faqCategories": faqCategoriesMenu(),
	"courses":       coursesMenu(),
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
func faqCategoriesMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("О школе"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Вопросы об оплате"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Вопросы об обучении"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад в главное меню"),
		),
	)
}

func faqPaymentMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Как и когда происходит оплата курса?"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Хочу оплатить в рассрочку. Какие условия?"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Я из другой страны. Могу ли я записаться на курс? Как проходит оплата?"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("назад"),
		),
	)
}

func faqStudyMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Как проходит обучение?"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Возможно ли взять перерыв во время курса?"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Выдается ли сертификат по окончании обучения?"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("назад"),
		),
	)
}

func faqAboutSchoolMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Подробнее про школу"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Как понять на какой курс я могу записаться со своим уровнем?"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Как я могу записаться на курс?"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("назад"),
		),
	)
}

//func faqMenu() tgbotapi.ReplyKeyboardMarkup {
//	return tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(
//			tgbotapi.NewKeyboardButton("Подробнее про школу"),
//		),
//		tgbotapi.NewKeyboardButtonRow(
//			tgbotapi.NewKeyboardButton("Как проходит обучение?"),
//		),
//		tgbotapi.NewKeyboardButtonRow(
//			tgbotapi.NewKeyboardButton("Как я могу записаться на курс?"),
//		),
//		//tgbotapi.NewKeyboardButtonRow(
//		//	tgbotapi.NewKeyboardButton("Когда стартует ближайший набор на курсы?"),
//		//),
//		//tgbotapi.NewKeyboardButtonRow(
//		//	tgbotapi.NewKeyboardButton("Формат обучения?"),
//		//),
//		tgbotapi.NewKeyboardButtonRow(
//			tgbotapi.NewKeyboardButton("Как и когда происходит оплата курса?"),
//		),
//		tgbotapi.NewKeyboardButtonRow(
//			tgbotapi.NewKeyboardButton("Хочу оплатить в рассрочку. Какие условия?"),
//		),
//		tgbotapi.NewKeyboardButtonRow(
//			tgbotapi.NewKeyboardButton("Я из другой страны. Могу ли я записаться на курс? Как проходит оплата?"),
//		),
//		tgbotapi.NewKeyboardButtonRow(
//			tgbotapi.NewKeyboardButton("Как понять на какой курс я могу записаться со своим уровнем?"),
//		),
//		tgbotapi.NewKeyboardButtonRow(
//			tgbotapi.NewKeyboardButton("Возможно ли взять перерыв во время курса?"),
//		),
//		tgbotapi.NewKeyboardButtonRow(
//			tgbotapi.NewKeyboardButton("Выдается ли сертификат по окончании обучения?"),
//		),
//		tgbotapi.NewKeyboardButtonRow(
//			tgbotapi.NewKeyboardButton("Назад в главное меню"),
//		),
//	)
//}

// courses

func coursesMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Основы рисунка")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Форма и тон")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Свет и цвет")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Портрет: Скетчинг и стилизация")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Скетчинг: тело, движение, одежда")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Динамический портрет")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Фигура человека")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Мастерская с Евой")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Дизайн существ")),
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

	if info.MainInfo != "" {
		buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Основная информация"),
		))
	}
	//if info.Duration != "" {
	//	buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
	//		tgbotapi.NewKeyboardButton("Длительность курса"),
	//	))
	//}

	//if info.StartDate != "" {
	//	buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
	//		tgbotapi.NewKeyboardButton("Ближайший старт"),
	//	))
	//}

	if info.Tariffs != "" {
		buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Доступные тарифы"),
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
	if info.ForWhom != "" {
		buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Для кого подходит курс"),
		))
	}

	if info.WhereToFindWorks != "" {
		buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Где посмотреть работы куратора и студентов?"),
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
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Основы рисунка")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Форма и тон")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Свет и цвет")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Портрет: Скетчинг и стилизация")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Скетчинг: тело, движение, одежда")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Динамический портрет")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Фигура человека")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Мастерская с Евой")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Дизайн существ")),
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

// ===== INLINE MENUS =====

func inlineButton(text, data string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(text, data)
}

func inlineRow(buttons ...tgbotapi.InlineKeyboardButton) []tgbotapi.InlineKeyboardButton {
	return buttons
}

// Главное меню
func mainInlineMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		inlineRow(
			inlineButton("Частые вопросы", "faq"),
		),
		inlineRow(
			inlineButton("Все курсы", "courses"),
		),
		inlineRow(
			inlineButton("Записаться в лист ожидания", "waitlist"),
		),
	)
}

// FAQ — категории
func faqCategoriesInlineMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		inlineRow(
			inlineButton("О школе", "faq_about"),
		),
		inlineRow(
			inlineButton("Вопросы об оплате", "faq_payment"),
		),
		inlineRow(
			inlineButton("Вопросы об обучении", "faq_study"),
		),
		inlineRow(
			inlineButton("Назад в главное меню", "main"),
		),
	)
}

// FAQ — об оплате
func faqPaymentInlineMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		inlineRow(
			inlineButton("Как и когда происходит оплата курса?", "faq_pay_when"),
		),
		inlineRow(
			inlineButton("Хочу оплатить в рассрочку. Какие условия?", "faq_installment"),
		),
		inlineRow(
			inlineButton("Я из другой страны. Могу ли я записаться на курс? Как проходит оплата?", "faq_foreign"),
		),
		inlineRow(
			inlineButton("назад", "faq"),
		),
	)
}

// FAQ — обучение
func faqStudyInlineMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		inlineRow(
			inlineButton("Как проходит обучение?", "faq_format"),
		),
		inlineRow(
			inlineButton("Возможно ли взять перерыв во время курса?", "faq_pause"),
		),
		inlineRow(
			inlineButton("Выдается ли сертификат по окончании курса?", "faq_certificate"),
		),
		inlineRow(
			inlineButton("назад", "faq"),
		),
	)
}

// FAQ — о школе
func faqAboutSchoolInlineMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		inlineRow(
			inlineButton("Подробнее про школу", "faq_about_info"),
		),
		inlineRow(
			inlineButton("Как понять на какой курс я могу записаться со своим уровнем?", "faq_level"),
		),
		inlineRow(
			inlineButton("Как я могу записаться на курс?", "faq_register"),
		),
		inlineRow(
			inlineButton("назад", "faq"),
		),
	)
}

// Курсы
func coursesInlineMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		inlineRow(
			inlineButton("Основы рисунка", "course:Основы рисунка"),
		),
		inlineRow(
			inlineButton("Форма и тон", "course:Форма и тон"),
		),
		inlineRow(
			inlineButton("Свет и цвет", "course:Свет и цвет"),
		),
		inlineRow(
			inlineButton("Портрет: Скетчинг и стилизация", "course:Портрет: Скетчинг и стилизация"),
		),
		inlineRow(
			inlineButton("Скетчинг: тело, движение, одежда", "course:Скетчинг: тело, движение, одежда"),
		),
		inlineRow(
			inlineButton("Динамический портрет", "course:Динамический портрет"),
		),
		inlineRow(
			inlineButton("Фигура человека", "course:Фигура человека"),
		),
		inlineRow(
			inlineButton("Мастерская с Евой", "course:Мастерская с Евой"),
		),
		inlineRow(
			inlineButton("Дизайн существ", "course:Дизайн существ"),
		),
		inlineRow(
			inlineButton("Назад в главное меню", "main"),
		),
	)
}

// Меню конкретного курса
func courseInlineMenu(course string) tgbotapi.InlineKeyboardMarkup {
	rows := [][]tgbotapi.InlineKeyboardButton{}

	info := CoursesInfo[course]

	if info.MainInfo != "" {
		rows = append(rows, inlineRow(
			inlineButton("Основная информация", "courseinfo:main"),
		))
	}

	if info.Tariffs != "" {
		rows = append(rows, inlineRow(
			inlineButton("Доступные тарифы", "courseinfo:tariffs"),
		))
	}

	if info.Schedule != "" {
		rows = append(rows, inlineRow(
			inlineButton("Программа курса", "courseinfo:schedule"),
		))
	}

	if info.About != "" {
		rows = append(rows, inlineRow(
			inlineButton("О чем курс", "courseinfo:about"),
		))
	}

	if info.Tools != "" {
		rows = append(rows, inlineRow(
			inlineButton("Что понадобится", "courseinfo:tools"),
		))
	}

	if info.ForWhom != "" {
		rows = append(rows, inlineRow(
			inlineButton("Для кого подходит курс", "courseinfo:forwhom"),
		))
	}

	if info.WhereToFindWorks != "" {
		rows = append(rows, inlineRow(
			inlineButton("Где посмотреть работы куратора и студентов?", "courseinfo:works"),
		))
	}

	rows = append(rows, inlineRow(
		inlineButton("Назад к списку курсов", "courses"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// Лист ожидания
func waitlistInlineMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		inlineRow(
			inlineButton("Основы рисунка", "waitlist:Основы рисунка"),
		),
		inlineRow(
			inlineButton("Форма и тон", "waitlist:Форма и тон"),
		),
		inlineRow(
			inlineButton("Свет и цвет", "waitlist:Свет и цвет"),
		),
		inlineRow(
			inlineButton("Портрет: Скетчинг и стилизация", "waitlist:Портрет:Скетчинг и стилизация"),
		),
		inlineRow(
			inlineButton("Скетчинг: тело, движение, одежда", "waitlist:Скетчинг: тело, движение, одежда"),
		),
		inlineRow(
			inlineButton("Динамический портрет", "waitlist:Динамический портрет"),
		),
		inlineRow(
			inlineButton("Фигура человека", "waitlist:Фигура человека"),
		),
		inlineRow(
			inlineButton("Мастерская с Евой", "waitlist:Мастерская с Евой"),
		),
		inlineRow(
			inlineButton("Дизайн существ", "waitlist:Дизайн существ"),
		),
		inlineRow(
			inlineButton("Назад в главное меню", "main"),
		),
	)
}
