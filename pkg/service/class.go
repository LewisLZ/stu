package service

import (
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/quexer/utee"

	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/ut"
	"liuyu/stu/pkg/web/hdl/form"
)

type Class struct {
	Ds *datasource.Ds
}

func (p *Class) Get(id int) (*model.Class, error) {
	var cla model.Class

	if err := p.Ds.Db.Table("class c").
		Select("c.id, c.school_year_id, c.name, s.year, s.pos").
		Joins("LEFT JOIN school_year s ON c.school_year_id=s.id").
		Where("c.id=?", id).Scan(&cla).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	tmp, err := time.ParseInLocation("2006", cla.Year, time.Local)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	cla.YearTmp = utee.Tick(tmp)

	return &cla, nil
}

func (p *Class) List(in *form.ListClass) ([]*model.Class, int, error) {

	db := p.Ds.Db.Table("class c").Select("c.id, c.school_year_id, c.name, s.year, s.pos").
		Joins("LEFT JOIN school_year s ON c.school_year_id=s.id")
	if in.Name != "" {
		db = db.Where("c.name like ?", "%"+in.Name+"%")
	}
	if in.Year != "" {
		db = db.Where("s.year = ?", in.Year)
	}
	if in.Pos >= model.Pos_Up && in.Pos <= model.Pos_Down {
		db = db.Where("s.pos = ?", in.Pos)
	}

	var count int
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	var clasies []*model.Class

	_, limit, offset := ut.MakePager(in.Page, in.Limit, 10)
	if err := db.Limit(limit).Offset(offset).Order("c.id desc").
		Scan(&clasies).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	for _, class := range clasies {
		teachers := []*model.Teacher{}
		if err := p.Ds.Db.Table(model.C_Teacher+" t").
			Joins("LEFT JOIN teacher_class tc ON tc.teacher_id=t.id").
			Select("t.id, t.name, t.mobile, t.sex").
			Where("tc.class_id=?", class.Id).
			Order("t.id desc").Scan(&teachers).Error; err != nil {
			return nil, 0, errors.WithStack(err)
		}
		class.Teacher = teachers

		for _, t := range class.Teacher {
			curriculum := []*model.Curriculum{}
			if err := p.Ds.Db.Table(model.C_Curriculum+" c").
				Select("c.id, c.name").
				Joins("LEFT JOIN teacher_curriculum tc ON c.id=tc.curriculum_id").
				Where("tc.teacher_id=?", t.Id).
				Scan(&curriculum).Error; err != nil {
				return nil, 0, errors.WithStack(err)
			}
			t.Curriculum = curriculum
		}

		students := []*model.Student{}
		if err := p.Ds.Db.Table(model.C_Student).
			Select("id, code, name, mobile, sex, birthday, intake_time, address").
			Where("class_id=?", class.Id).
			Order("id desc").Scan(&students).Error; err != nil {
			return nil, 0, errors.WithStack(err)
		}
		class.Student = students
	}

	return clasies, count, nil
}

func (p *Class) ListNameByIds(ids []int) ([]string, error) {
	var names []string

	var res []*struct {
		Name string
		Year string
		Pos  model.Pos
	}

	if err := p.Ds.Db.Raw(`select c.name, s.year, s.pos
			from class c, school_year s 
			where c.id in (?) and c.school_year_id=s.id`, ids).Order("c.id desc").Scan(&res).
		Error; err != nil {
		return nil, err
	}

	for _, v := range res {
		names = append(names, ut.MakeClassName(v.Name, v.Year, v.Pos))
	}

	return names, nil
}

func (p *Class) Save(in *form.SaveClass) error {
	class := model.Class{}
	tick := utee.Tick()
	if err := copier.Copy(&class, in); err != nil {
		return err
	}
	if in.Id > 0 {
		if err := p.Ds.Db.First(&class, in.Id).Error; err != nil {
			return err
		}
	} else {
		class.Ct = tick
	}
	class.Mt = tick
	class.Name = in.Name

	if err := p.Ds.Db.Save(&class).Error; err != nil {
		return err
	}
	return nil
}
