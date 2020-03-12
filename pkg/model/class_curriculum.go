package model

// 课程
type ClassCurriculum struct {
	Base
	ClassId      int `gorm:"not null;unique_index:u_idx_classid_cid_classcurr"`
	CurriculumId int `gorm:"not null;unique_index:u_idx_classid_cid_classcurr"`
}
