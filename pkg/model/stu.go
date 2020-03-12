package model

// 学生
type Stu struct {
	Base
	ClassId    int    `gorm:"not null"`
	Name       string `gorm:"not null;size:10"`
	Code       string `gorm:"not null;size:10"`
	Sex        int    `gorm:"not null"` // 性别 1:男，2:女
	Birthday   int64  `gorm:"not null;default:0"`
	Address    string `gorm:"not null;size:200"`
	IntakeTime int64  `gorm:"not null;default:0"`
	Mobile     string `gorm:"not null;size:11"`
}
