package model

const (
	C_User = "user"
)

// User 用户
type User struct {
	Base
	Account string `gorm:"not null;size:80"` // 登录名
	Passwd  string `gorm:"not null;size:32"` // 密码
	Name    string `gorm:"not null;size:20"` // 姓名
	Mobile  string `gorm:"not null;size:11"` // 手机号
}
