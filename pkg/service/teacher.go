package service

import (
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"github.com/quexer/utee"

	"liuyu/stu/pkg/dao"
	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/ut"
	"liuyu/stu/pkg/web/hdl/form"
)

type Teacher struct {
	Ds            *datasource.Ds
	ClassDao      *dao.Class
	CurriculumDao *dao.Curriculum
}

func (p *Teacher) Get(id int) (*model.Teacher, error) {
	var teacher model.Teacher

	if err := p.Ds.Db.Model(model.Teacher{}).Select("id, name, sex, mobile").Where("id=?", id).Scan(&teacher).Error; err != nil {
		return nil, err
	}

	if err := p.MakeClassAndCurriculumIds(&teacher); err != nil {
		return nil, err
	}

	classIds := []int{}
	for _, v := range teacher.Class {
		classIds = append(classIds, v.Id)
	}
	teacher.ClassIds = classIds

	curriculumIds := []int{}
	for _, v := range teacher.Curriculum {
		curriculumIds = append(curriculumIds, v.Id)
	}
	teacher.CurriculumIds = curriculumIds

	return &teacher, nil
}

func (p *Teacher) List(in *form.ListTeacher) ([]*model.Teacher, int, error) {
	_, limit, offset := ut.MakePager(in.Page, in.Limit, 10)

	db := p.Ds.Db.Model(&model.Teacher{}).Select("id, name, sex, mobile")
	if in.Name != "" {
		db = db.Where("name like ?", "%"+in.Name+"%")
	}
	if in.Mobile != "" {
		db = db.Where("mobile like ?", "%"+in.Mobile+"%")
	}
	if in.Sex > model.SexUnknown {
		db = db.Where("sex = ?", in.Sex)
	}

	var teachers []*model.Teacher
	var count int
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Limit(limit).Offset(offset).Order("id desc").Scan(&teachers).Error; err != nil {
		return nil, 0, err
	}

	for _, t := range teachers {
		if err := p.MakeClassAndCurriculumIds(t); err != nil {
			return nil, 0, err
		}
	}
	return teachers, count, nil
}

func (p *Teacher) MakeClassAndCurriculumIds(t *model.Teacher) error {
	class, err := p.ClassDao.ListClassWithTeacherId(p.Ds.Db, t.Id)
	if err != nil {
		return err
	}
	t.Class = class

	curriculum, err := p.CurriculumDao.ListCurriculumWithTeacherId(p.Ds.Db, t.Id)
	if err != nil {
		return err
	}
	t.Curriculum = curriculum
	return nil
}

func (p *Teacher) Save(in *form.SaveTeacher) error {
	teacher := model.Teacher{}
	tick := utee.Tick()
	if in.Id > 0 {
		if err := p.Ds.Db.First(&teacher, in.Id).Error; err != nil {
			return err
		}
	} else {
		teacher.Ct = tick
	}
	if err := copier.Copy(&teacher, in); err != nil {
		return err
	}
	teacher.Mt = tick

	err := datasource.RunTransaction(p.Ds.Db, func(tx *gorm.DB) error {
		if err := tx.Save(&teacher).Error; err != nil {
			return err
		}

		if err := tx.Where("teacher_id=?", teacher.Id).Delete(&model.TeacherClass{}).Error; err != nil {
			return err
		}

		if err := tx.Where("teacher_id=?", teacher.Id).Delete(&model.TeacherCurriculum{}).Error; err != nil {
			return err
		}

		for _, id := range in.ClassIds {
			tick := utee.Tick()
			tc := model.TeacherClass{
				Base: model.Base{
					Ct: tick,
					Mt: tick,
				},
				TeacherId: teacher.Id,
				ClassId:   id,
			}
			if err := tx.Model(&model.TeacherClass{}).Create(&tc).Error; err != nil {
				return err
			}
		}

		for _, id := range in.CurriculumIds {
			tick := utee.Tick()
			tc := model.TeacherCurriculum{
				Base: model.Base{
					Ct: tick,
					Mt: tick,
				},
				TeacherId:    teacher.Id,
				CurriculumId: id,
			}
			if err := tx.Model(&model.TeacherCurriculum{}).Create(&tc).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
