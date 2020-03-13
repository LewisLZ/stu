package model

const (
	C_Class = "class"
)

// 班级
type Class struct {
	Base
	ParentId int    `gorm:"not null;unique_index:uidx_parent_name_class" json:"parent_id"` // 上级Id
	Name     string `gorm:"not null;unique_index:uidx_parent_name_class" json:"name"`      // 班级

	Children []*Class   `gorm:"-" json:"children"` // 子结构
	Teacher  []*Teacher `gorm:"-" json:"teacher"`
	Student  []*Student `gorm:"-" json:"student"`
}
