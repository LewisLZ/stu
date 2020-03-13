package model

// 学生
type Student struct {
	Base
	ClassId    int    `gorm:"not null" json:"class_id"`           // 班级Id
	Name       string `gorm:"not null;size:10" json:"name"`       // 姓名
	Code       string `gorm:"not null;size:15" json:"code"`       // 学号
	Sex        Sex    `gorm:"not null" json:"sex"`                // 性别 1:男，2:女
	Birthday   string `gorm:"not null;size:10" json:"birthday"`   // 生日
	Address    string `gorm:"not null;size:200" json:"address"`   // 地址
	IntakeTime string `gorm:"not null;size:7" json:"intake_time"` // 入学时间
	Mobile     string `gorm:"not null;size:11" json:"mobile"`     // 联系方式

	ClassName     string `gorm:"-" json:"class_name"`
	BirthdayTmp   int64  `gorm:"-" json:"birthday_tmp"`    // 生日
	IntakeTimeTmp int64  `gorm:"-" json:"intake_time_tmp"` // 入学时间
}

type StudentCode struct {
	Base
	IntakeTime string `gorm:"not null;size:7;unique_index:uidx_intake_code_scode"`  // 入学时间
	Code       string `gorm:"not null;size:15;unique_index:uidx_intake_code_scode"` // 学号
}
