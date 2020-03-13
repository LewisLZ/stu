package model

// 学生
type Stu struct {
	Base
	ClassId    int    `gorm:"not null"`          // 班级Id
	Name       string `gorm:"not null;size:10"`  // 姓名
	Code       string `gorm:"not null;size:10"`  // 学号
	Sex        Sex    `gorm:"not null"`          // 性别 1:男，2:女
	Birthday   string `gorm:"not null;size:10"`  // 生日
	Address    string `gorm:"not null;size:200"` // 地址
	IntakeTime string `gorm:"not null;size:10"`  // 入学时间
	Mobile     string `gorm:"not null;size:11"`  // 联系方式
}
