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
