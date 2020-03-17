package model

const (
	C_SchoolYear = "school_year"
	C_Class      = "class"
)

type Pos int

const (
	Pos_Up   Pos = 1
	Pos_Down Pos = 2
)

type SchoolYear struct {
	Base
	Year string `gorm:"not null;size:4;unique_index:uidx_year_pos_sy" json:"year"`
	Pos  Pos    `gorm:"not null;size:1;unique_index:uidx_year_pos_sy" json:"pos"`

	YearTmp int64    `gorm:"-" json:"year_tmp"`
	Class   []*Class `gorm:"-" json:"class"`
}

// 班级
type Class struct {
	Base
	SchoolYearId int    `gorm:"not null;unique_index:uidx_parent_name_class" json:"school_year_id"` // 上级Id
	Name         string `gorm:"not null;size:4;unique_index:uidx_parent_name_class" json:"name"`    // 班级

	Teacher                  []*Teacher `gorm:"-" json:"teacher"`
	Student                  []*Student `gorm:"-" json:"student"`
	Year                     string     `gorm:"-" json:"year"`
	Pos                      Pos        `gorm:"-" json:"pos"`
	YearTmp                  int64      `gorm:"-" json:"year_tmp"`
	ClassCurriculumYearCount int        `gorm:"-" json:"class_curriculum_year_count"`
}
