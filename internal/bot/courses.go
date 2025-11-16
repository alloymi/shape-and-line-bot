package bot

type CourseInfo struct {
	Name      string
	Duration  string
	StartDate string
	Curator   string
}

var CoursesInfo = map[string]CourseInfo{
	"Фигура человека": {
		Name:      "Фигура человека",
		Duration:  "1 недель",
		StartDate: "1",
		Curator:   "Куратор: 1",
	},

	"": {
		Name:      "",
		Duration:  " недель",
		StartDate: "",
		Curator:   "Куратор: ",
	},
}
