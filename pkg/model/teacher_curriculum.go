package model

type TeacherCurriculum struct {
	Base
	TeacherId    int `gorm:"not null;unique_index:u_idx_tid_cid_tcu"`
	CurriculumId int `gorm:"not null;unique_index:u_idx_tid_cid_tcu"`
}
