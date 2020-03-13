package form

import "liuyu/stu/pkg/model"

type SaveSchoolYear struct {
	Id   int       `json:"id"`
	Year string    `json:"year"`
	Pos  model.Pos `json:"pos"`
}

type ListSchoolYear struct {
	Page  int       `form:"page"`
	Limit int       `form:"limit"`
	Year  string    `form:"year"`
	Pos   model.Pos `form:"pos"`
}
