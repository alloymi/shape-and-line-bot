package bot

type BotState string

const (
	StateDefault BotState = "default"
	StateFAQ     BotState = "faq"
	StateCourses BotState = "courses"
)

var userState = map[int64]BotState{}

func SetState(chatID int64, s BotState) {
	userState[chatID] = s
}

func GetState(chatID int64) BotState {
	if s, ok := userState[chatID]; ok {
		return s
	}
	return StateDefault
}
