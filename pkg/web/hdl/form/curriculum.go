package form

type SaveCurriculum struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ListCurriculum struct {
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
	Name  string `form:"name"`
}

type ListCurriculumForExamination struct {
	ClassId int `form:"class_id"`
}

type ListCurriculumChoose struct {
	CCYearId int `form:"cc_year_id"`
}
