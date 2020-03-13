package model

const (
	C_Teacher = "teacher"
)

type Sex int

const (
	SexUnknown Sex = 0
	SexBoy     Sex = 1
	SexGirl    Sex = 2
)

// 教师
type Teacher struct {
	Base
	Name   string `gorm:"not null;size:10" json:"name,omitempty"`
	Sex    Sex    `gorm:"not null" json:"sex,omitempty"` // 性别 1:男，2:女
	Mobile string `gorm:"not null;size:11" json:"mobile,omitempty"`

	Class         []*Class      `gorm:"-" json:"class"`
	Curriculum    []*Curriculum `gorm:"-" json:"curriculum"`
	ClassIds      []int         `gorm:"-" json:"class_ids"`
	CurriculumIds []int         `gorm:"-" json:"curriculum_ids"`
}
