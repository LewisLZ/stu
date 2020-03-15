package model

type Achievement struct {
	Base
	ExaminationId         int    `json:"examination_id"`
	ExaminationName       string `json:"examination_name"`
	ExaminationTime       int64  `json:"examination_time"`
	ExaminationClassCount int    `json:"examination_class_count"`
}
