package model

const (
	C_User = "user"
)

type UserType int

const (
	UserType_Normal = 0
	UserType_Super  = 1
)

// User 用户
type User struct {
	Base
	Mobile string   `gorm:"not null;size:11;unique" json:"mobile"` // 手机号
	Passwd string   `gorm:"not null;size:32" json:"passwd"`        // 密码
	Name   string   `gorm:"not null;size:20" json:"name"`          // 姓名
	Type   UserType `gorm:"not null;size:2" json:"type"`           // 是否是超级管理员
}
