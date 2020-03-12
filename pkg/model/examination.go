package model

type Examination struct {
	Base
	Name   string `gorm:"not null;size:10"`
	Remark string `gorm:"not null;size:200"`
}

type ExaminationItem struct {
	Base
	ExaminationId int   `gorm:"not null"`
	Time          int64 `gorm:"not null:default:0"`
	CurriculumId  int   `gorm:"not null"`
	ClassId       int   `gorm:"not null"`
}
