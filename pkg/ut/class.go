package ut

import (
	"fmt"

	"liuyu/stu/pkg/model"
)

func MakeClassName(className, year string, pos model.Pos) string {
	p := "上半年"
	if pos == model.Pos_Down {
		p = "下半年"
	}
	return fmt.Sprintf("%s(%s-%s)", className, year, p)
}
