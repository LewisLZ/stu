package model

type TeacherClass struct {
	Base
	TeacherId int `gorm:"not null;unique_index:u_idx_tid_cid_tc"`
	ClassId   int `gorm:"not null;unique_index:u_idx_tid_cid_tc"`
	Admin     int `gorm:"not null;default:0"` // 1 表示班主任
}
