package model

type ClassCurriculumYear struct {
	Base
	ClassId int    `gorm:"not null;unique_index:uidx_cyp_ccy" json:"class_id"`
	Year    string `gorm:"not null;size:4;unique_index:uidx_cyp_ccy" json:"year"`
	Pos     Pos    `gorm:"not null;size:1;unique_index:uidx_cyp_ccy" json:"pos"`

	ClassName       string             `gorm:"-" json:"class_name"`
	ClassCurriculum []*ClassCurriculum `gorm:"-" json:"class_curriculum"`
	YearTmp         int64              `gorm:"-" json:"year_tmp"`
	ClassYear       string             `gorm:"-" json:"class_year"`
	ClassYearTmp    int64              `gorm:"-" json:"class_year_tmp"`
	ClassPos        Pos                `gorm:"-" json:"class_pos"`
}

// 课程
type ClassCurriculum struct {
	Base
	CCYearId     int `gorm:"not null;unique_index:u_idx_ccyearid_cid_classcurr" json:"cc_year_id"`
	CurriculumId int `gorm:"not null;unique_index:u_idx_ccyearid_cid_classcurr" json:"curriculum_id"`

	CurriculumName string `gorm:"-" json:"curriculum_name"`
}
