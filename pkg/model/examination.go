package model

type Examination struct {
	Base
	Name      string `gorm:"not null;size:10" json:"name"`
	StartTime int64  `gorm:"not null" json:"start_time"`
	Remark    string `gorm:"not null;size:200;default:''" json:"remark"`

	ExaminationItemCount int `gorm:"-" json:"examination_item_count"`
}

type ExaminationClass struct {
	Base
	ExaminationId int `gorm:"not null" json:"examination_id"`
	ClassId       int `gorm:"not null" json:"class_id"`
}

type ExaminationCurriculum struct {
	Base
	ExaminationClassId int `gorm:"not null" json:"examination_class_id"`
	CurriculumId       int `gorm:"not null" json:"curriculum_id"`
}
