package form

import "liuyu/stu/pkg/model"

type SaveStudent struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	ClassId    int       `json:"class_id"`
	Sex        model.Sex `json:"sex"`
	Birthday   string    `json:"birthday"`
	Address    string    `json:"address"`
	IntakeTime string    `json:"intake_time"`
	Mobile     string    `json:"mobile"`
}

type ListStudent struct {
	Page       int       `form:"page"`
	Limit      int       `form:"limit"`
	Code       string    `form:"code"`
	Sex        model.Sex `json:"sex"`
	Name       string    `form:"name"`
	Birthday   string    `form:"birthday"`
	Address    string    `form:"address"`
	IntakeTime string    `form:"intake_time"`
	Mobile     string    `form:"mobile"`
}
