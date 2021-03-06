package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/quexer/utee"

	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/ut"
	"liuyu/stu/pkg/web/hdl/form"
)

type Student struct {
	Ds *datasource.Ds
}

func (p *Student) Get(id int) (*model.Student, error) {
	var stu model.Student

	if err := p.Ds.Db.Table("student t, class c, school_year s").
		Select("t.id, t.name, t.code, t.sex, t.class_id, t.mobile, t.birthday, t.intake_time, t.address, c.name class_name, s.year, s.pos").
		Where("t.class_id=c.id AND c.school_year_id=s.id AND t.id=?", id).Scan(&stu).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	stu.ClassName = ut.MakeClassName(stu.ClassName, stu.Year, stu.Pos)

	return &stu, nil
}

func (p *Student) List(in *form.ListStudent) ([]*model.Student, int, error) {
	db := p.Ds.Db.Table("student t, class c, school_year s").
		Select("t.id, t.name, t.code, t.sex, t.mobile, t.birthday, t.intake_time, t.address, c.name class_name, s.year, s.pos").
		Where("t.class_id=c.id AND c.school_year_id=s.id")

	if in.IntakeTime != "" {
		db = db.Where("t.intake_time=?", in.IntakeTime)
	}
	if in.Code != "" {
		db = db.Where("t.code like ?", "%"+in.Code+"%")
	}
	if in.Sex > model.SexUnknown && in.Sex <= model.SexGirl {
		db = db.Where("t.sex = ?", in.Sex)
	}
	if in.Name != "" {
		db = db.Where("t.name like ?", "%"+in.Name+"%")
	}
	if in.Birthday != "" {
		db = db.Where("t.birthday=?", in.Birthday)
	}
	if in.Address != "" {
		db = db.Where("t.address like ?", "%"+in.Address+"%")
	}
	if in.Mobile != "" {
		db = db.Where("t.mobile like ?", "%"+in.Mobile+"%")
	}

	var count int
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	_, limit, offset := ut.MakePager(in.Page, in.Limit, 10)
	var stus []*model.Student
	if err := db.Limit(limit).Offset(offset).Order("t.id desc").Scan(&stus).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	for _, stu := range stus {
		stu.ClassName = ut.MakeClassName(stu.ClassName, stu.Year, stu.Pos)
	}
	return stus, count, nil
}

func (p *Student) Save(in *form.SaveStudent) error {
	stu := model.Student{}
	tick := utee.Tick()
	if in.Id > 0 {
		if err := p.Ds.Db.First(&stu, in.Id).Error; err != nil {
			return errors.WithStack(err)
		}
	} else {
		code, err := p.makeCode(in)
		if err != nil {
			return errors.WithStack(err)
		}
		stu.Ct = tick
		stu.Code = code
	}
	if err := copier.Copy(&stu, in); err != nil {
		return errors.WithStack(err)
	}
	stu.Mt = tick

	err := datasource.RunTransaction(p.Ds.Db, func(tx *gorm.DB) error {
		if err := tx.Save(&stu).Error; err != nil {
			return errors.WithStack(err)
		}
		if in.Id == 0 {
			scode := model.StudentCode{Code: stu.Code, IntakeTime: in.IntakeTime}
			scode.Ct = tick
			scode.Mt = tick
			if err := tx.Create(&scode).Error; err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (p *Student) makeCode(stu *form.SaveStudent) (string, error) {
	intakeTime := strings.Replace(stu.IntakeTime, "-", "", -1)

	var res struct {
		Code string
	}
	err := p.Ds.Db.Model(&model.StudentCode{}).
		Select("code").
		Where("intake_time=?", stu.IntakeTime).
		Order("id desc").Limit(1).Scan(&res).Error
	if !datasource.NotFound(err) && err != nil {
		return "", errors.WithStack(err)
	}
	var code = 0
	if res.Code != "" {
		c := strings.Replace(res.Code, intakeTime, "", -1)
		code, err = strconv.Atoi(c)
		if err != nil {
			return "", errors.WithStack(err)
		}
	}

	return fmt.Sprintf("%s%05d", intakeTime, code+1), nil
}
