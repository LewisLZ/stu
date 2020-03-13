package dao

import (
	"github.com/jinzhu/gorm"

	"liuyu/stu/pkg/model"
)

type Curriculum struct {
}

func (p *Curriculum) ListCurriculumWithTeacherId(db *gorm.DB, tid int) ([]*model.Curriculum, error) {
	curriculum := []*model.Curriculum{}
	if err := db.Table(model.C_Curriculum+" c").
		Select("c.id, c.name").
		Joins("LEFT JOIN teacher_curriculum tc ON c.id=tc.curriculum_id").
		Where("tc.teacher_id=?", tid).
		Scan(&curriculum).Error; err != nil {
		return nil, err
	}

	return curriculum, nil
}
