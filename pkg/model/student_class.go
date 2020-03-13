package model

type StudentClass struct {
	Base
	StuId   int `gorm:"not null;unique_index:u_idx_sid_cid_tcl"`
	ClassId int `gorm:"not null;unique_index:u_idx_sid_cid_tcl"`
}
