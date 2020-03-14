package model

const (
	C_Curriculum = "curriculum"
)

// 课程
type Curriculum struct {
	Base
	Name string `gorm:"not null;size:20" json:"name,omitempty"`

	Teacher           []*Teacher `gorm:"-" json:"teacher"`
	Year              string     `gorm:"-" json:"-"`
	Pos               Pos        `gorm:"-" json:"-"`
	ClassCurriculumId int        `gorm:"-" json:"class_curriculum_id"`
}
