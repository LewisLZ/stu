package form

import "liuyu/stu/pkg/model"

type SaveClass struct {
	Id           int    `json:"id"`
	SchoolYearId int    `json:"school_year_id"`
	Name         string `json:"name"`
}

type ListClass struct {
	Page  int       `form:"page"`
	Limit int       `form:"limit"`
	Name  string    `form:"name"`
	Year  string    `form:"year"`
	Pos   model.Pos `form:"pos"`
}
