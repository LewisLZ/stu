package form

import "liuyu/stu/pkg/model"

type SaveClassCurriculumYear struct {
	Id      int       `json:"id"`
	ClassId int       `json:"class_id"`
	Year    string    `json:"year"`
	Pos     model.Pos `json:"pos"`
}

type ListClassCurriculumYear struct {
	ClassId int `form:"class_id"`
}

type SaveClassCurriculum struct {
	CCYearId      int   `json:"cc_year_id"`
	CurriculumIds []int `json:"curriculum_ids"`
}
