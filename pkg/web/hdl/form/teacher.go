package form

type SaveTeacher struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Sex           int    `json:"sex"`
	Mobile        string `json:"mobile"`
	CurriculumIds []int  `json:"curriculum_ids"`
	ClassIds      []int  `json:"class_ids"`
}

type ListTeacher struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Name   string `form:"name"`
	Mobile string `form:"mobile"`
}
