package form

import "liuyu/stu/pkg/model"

type SaveTeacher struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Sex           int    `json:"sex"`
	Mobile        string `json:"mobile"`
	CurriculumIds []int  `json:"curriculum_ids"`
	ClassIds      []int  `json:"class_ids"`
}

type ListTeacher struct {
	Page   int       `form:"page"`
	Limit  int       `form:"limit"`
	Sex    model.Sex `form:"sex"`
	Name   string    `form:"name"`
	Mobile string    `form:"mobile"`
}
