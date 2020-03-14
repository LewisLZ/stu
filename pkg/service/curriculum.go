package service

import (
	"github.com/jinzhu/copier"
	"github.com/quexer/utee"

	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/ut"
	"liuyu/stu/pkg/web/hdl/form"
)

type Curriculum struct {
	Ds *datasource.Ds
}

func (p *Curriculum) ListForExamination(req *form.ListCurriculumForExamination) ([]*model.Curriculum, error) {
	db := p.Ds.Db.Table("curriculum c").Select("c.id, c.name, ccy.year, ccy.pos, cc.id class_curriculum_id").
		Joins("LEFT JOIN class_curriculum cc ON cc.curriculum_id=c.id").
		Joins("LEFT JOIN class_curriculum_year ccy ON ccy.id=cc.cc_year_id").
		Joins("LEFT JOIN class cl ON ccy.class_id=cl.id").
		Where("cl.id=?", req.ClassId)

	var curriculums []*model.Curriculum
	if err := db.Order("c.id desc").Scan(&curriculums).Error; err != nil {
		return nil, err
	}
	for _, curr := range curriculums {
		curr.Name = ut.MakeClassName(curr.Name, curr.Year, curr.Pos)
	}

	return curriculums, nil
}

func (p *Curriculum) ListChoose(req *form.ListCurriculumChoose) ([]*model.Curriculum, []int, error) {
	var curriculums []*model.Curriculum
	if err := p.Ds.Db.Model(&model.Curriculum{}).
		Select("id, name").Order("id desc").
		Scan(&curriculums).Error; err != nil {
		return nil, nil, err
	}

	var disabledIds []int
	if err := p.Ds.Db.Table("class_curriculum cc").
		Joins("RIGHT JOIN examination_class_curriculum ecc ON ecc.class_curriculum_id=cc.id").
		Joins("LEFT JOIN curriculum c ON c.id=cc.curriculum_id").
		Select("c.id").
		Where("cc.cc_year_id=?", req.CCYearId).
		Pluck("c.id", &disabledIds).
		Error; err != nil {
		return nil, nil, err
	}

	return curriculums, disabledIds, nil
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
	if err := db.Limit(limit).Offset(offset).Order("id desc").Scan(&curriculums).Error; err != nil {
		return nil, 0, err
	}

	for _, v := range curriculums {
		teachers := []*model.Teacher{}
		if err := p.Ds.Db.Table(model.C_Teacher+" t").
			Joins("LEFT JOIN teacher_curriculum tc ON tc.teacher_id=t.id").
			Select("t.id, t.name, t.mobile").
			Where("tc.curriculum_id=?", v.Id).Order("t.id desc").Scan(&teachers).Error; err != nil {
			return nil, 0, err
		}
		v.Teacher = teachers
	}

	return curriculums, count, nil
}

func (p *Curriculum) ListNameByIds(ids []int) ([]string, error) {
	var names []string
	if err := p.Ds.Db.Model(&model.Curriculum{}).Select("name").Where("id in (?)", ids).Order("id").Pluck("name", &names).Error; err != nil {
		return nil, err
	}
	return names, nil
}

func (p *Curriculum) ListNameByClassCurriculumIds(ids []int) ([]string, error) {
	var res []*struct {
		Name string
		Year string
		Pos  model.Pos
	}
	if err := p.Ds.Db.Table("curriculum c").
		Joins("LEFT JOIN class_curriculum cc ON cc.curriculum_id=c.id").
		Joins("LEFT JOIN class_curriculum_year ccy ON ccy.id=cc.cc_year_id").
		Select("c.name, ccy.year, ccy.pos").
		Where("cc.id in (?)", ids).
		Order("c.id").
		Scan(&res).
		Error; err != nil {
		return nil, err
	}

	var names []string
	for _, r := range res {
		names = append(names, ut.MakeClassName(r.Name, r.Year, r.Pos))
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
	if err := copier.Copy(&curriculum, req); err != nil {
		return err
	}
	curriculum.Mt = pick

	if err := p.Ds.Db.Save(&curriculum).Error; err != nil {
		return err
	}
	return nil
}
