package bot

type BotState string

const (
	StateDefault    BotState = "default"
	StateFAQ        BotState = "faq"
	StateFAQPayment BotState = "faq_payment"
	StateFAQStudy   BotState = "faq_study"
	StateFAQAbout   BotState = "faq_about"
	StateCourses    BotState = "courses"
	StateCourseMenu BotState = "course_menu"

	StateWaitlistChooseCourse BotState = "waitlist_choose_course"
	StateWaitlistAskFullName  BotState = "waitlist_ask_full_name"
	StateWaitlistAskEmail     BotState = "waitlist_ask_email"
)

var userState = map[int64]BotState{}
var userTempCourse = map[int64]string{}
var userTempFullname = map[int64]string{}

func SetState(chatID int64, s BotState) {
	userState[chatID] = s
}

func GetState(chatID int64) BotState {
	if s, ok := userState[chatID]; ok {
		return s
	}
	return StateDefault
}

func ResetState(chatID int64) {
	delete(userState, chatID)
	delete(userTempCourse, chatID)
	delete(userTempFullname, chatID)
}
