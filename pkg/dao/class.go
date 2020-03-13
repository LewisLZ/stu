package dao

import (
	"github.com/jinzhu/gorm"

	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/ut"
)

type Class struct {
}

func (p *Class) ListClassWithTeacherId(db *gorm.DB, tid int) ([]*model.Class, error) {
	class := []*model.Class{}
	if err := db.Table("class c, school_year s, teacher_class tc").
		Select("c.id, s.year, s.pos, c.name").
		Where("tc.teacher_id=? AND c.school_year_id=s.id AND c.id=tc.class_id", tid).
		Scan(&class).Error; err != nil {
		return nil, err
	}

	for _, c := range class {
		c.Name = ut.MakeClassName(c.Name, c.Year, c.Pos)
	}

	return class, nil
}
