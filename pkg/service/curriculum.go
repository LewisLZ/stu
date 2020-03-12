package service

import (
	"github.com/quexer/utee"

	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/ut"
	"liuyu/stu/pkg/web/hdl/form"
)

type Curriculum struct {
	Ds *datasource.Ds
}

func (p *Curriculum) List(req *form.ListCurriculum) ([]*model.Curriculum, int, error) {
	db := p.Ds.Db.Model(&model.Curriculum{}).Select("id, name")
	if req.Name != "" {
		db = db.Where("name like ?", "%"+req.Name+"%")
	}
	var count int
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := ut.MakePager(req.Page, req.Limit, 10)
	var curriculums []*model.Curriculum
	if err := db.Limit(limit).Offset(offset).Scan(&curriculums).Error; err != nil {
		return nil, 0, err
	}

	for _, v := range curriculums {
		teachers := []*model.Teacher{}
		if err := p.Ds.Db.Table(model.C_Teacher+" t").
			Joins("LEFT JOIN teacher_curriculum tc ON tc.teacher_id=t.id").
			Select("t.id, t.name, t.mobile").
			Where("tc.curriculum_id=?", v.Id).Scan(&teachers).Error; err != nil {
			return nil, 0, err
		}
		v.Teacher = teachers
	}

	return curriculums, count, nil
}

func (p *Curriculum) ListNameByIds(ids []int) ([]string, error) {
	var names []string
	if err := p.Ds.Db.Model(&model.Curriculum{}).Select("name").Where("id in (?)", ids).Pluck("name", &names).Error; err != nil {
		return nil, err
	}
	return names, nil
}

func (p *Curriculum) Save(req *form.SaveCurriculum) error {
	var curriculum model.Curriculum
	pick := utee.Tick()
	if req.Id != 0 {
		if err := p.Ds.Db.First(&curriculum, req.Id).Error; err != nil {
			return err
		}
	} else {
		curriculum.Ct = pick
	}
	curriculum.Mt = pick
	curriculum.Name = req.Name

	if err := p.Ds.Db.Save(&curriculum).Error; err != nil {
		return err
	}
	return nil
}
