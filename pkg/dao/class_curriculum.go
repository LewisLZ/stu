package dao

import (
	"github.com/jinzhu/gorm"

	"liuyu/stu/pkg/model"
)

type ClassCurriculum struct {
}

func (p *ClassCurriculum) ListClassCurriculumWithYearId(db *gorm.DB, id int) ([]*model.ClassCurriculum, error) {
	var cc []*model.ClassCurriculum

	if err := db.Table("class_curriculum cc").
		Select("cc.id, cc.cc_year_id, cc.curriculum_id, c.name curriculum_name").
		Joins("LEFT JOIN curriculum c ON c.id=cc.curriculum_id").
		Where("cc.cc_year_id=?", id).
		Error; err != nil {
		return nil, err
	}

	return cc, nil
}
