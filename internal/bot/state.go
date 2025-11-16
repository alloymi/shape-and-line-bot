package bot

type BotState string

const (
	StateDefault BotState = "default"
	StateFAQ     BotState = "faq"
	StateCourses BotState = "courses"

	StateWaitlistChooseCourse BotState = "waitlist_choose_course"
	StateWaitlistAskFullName  BotState = "waitlist_ask_full_name"
	StateWaitlistAskEmail     BotState = "waitlist_ask_email"

	StateCourseDetails BotState = "course_details"
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
