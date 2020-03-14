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

	Year                       string                        `gorm:"-" json:"-"`
	Pos                        Pos                           `gorm:"-" json:"-"`
	ClassName                  string                        `gorm:"-" json:"class_name"`
	ExaminationClassCurriculum []*ExaminationClassCurriculum `gorm:"-" json:"examination_class_curriculum"`
}

type ExaminationClassCurriculum struct {
	Base
	ExaminationClassId int `gorm:"not null" json:"examination_class_id"`
	ClassCurriculumId  int `gorm:"not null" json:"class_curriculum_id"`

	ClassCurriculumName string `gorm:"-" json:"class_curriculum_name"`
	ClassCurriculumYear string `gorm:"-" json:"-"`
	ClassCurriculumPos  Pos    `gorm:"-" json:"-"`
}
