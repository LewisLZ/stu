package model

type Base struct {
	Id int   `gorm:"primary_key" json:"id,omitempty"`
	Ct int64 `gorm:"not null;default:0" json:"ct,omitempty"`
	Mt int64 `gorm:"not null;default:0" json:"mt,omitempty"`
	Dt int64 `gorm:"not null;default:0" json:"dt,omitempty"`
}
