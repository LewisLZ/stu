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
		Where("cc.cc_year_id=?", id).Scan(&cc).
		Error; err != nil {
		return nil, err
	}

	for _, c := range cc {
		var curriculum struct {
			Name string
		}
		if err := db.Table(model.C_Curriculum).
			Select("name").
			Where("id=?", c.CurriculumId).
			Scan(&curriculum).Error; err != nil {
			return nil, err
		}

		c.CurriculumName = curriculum.Name
	}

	return cc, nil
}
