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

type SchoolYear struct {
	Ds *datasource.Ds
}

func (p *SchoolYear) List(in *form.ListSchoolYear) ([]*model.SchoolYear, int, error) {

	db := p.Ds.Db.Model(&model.SchoolYear{}).Select("id, year, pos")
	if in.Year != "" {
		db = db.Where("year = ?", in.Year)
	}
	if in.Pos >= model.Pos_Up && in.Pos <= model.Pos_Down {
		db = db.Where("pos = ?", in.Pos)
	}

	var count int
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	_, limit, offset := ut.MakePager(in.Page, in.Limit, 10)

	var schoolYears []*model.SchoolYear

	if err := db.Limit(limit).Offset(offset).Order("id desc").Scan(&schoolYears).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	for _, sy := range schoolYears {
		tmp, err := time.ParseInLocation("2006", sy.Year, time.Local)
		if err != nil {
			return nil, 0, errors.WithStack(err)
		}
		sy.YearTmp = utee.Tick(tmp)

		class := []*model.Class{}
		if err := p.Ds.Db.Table(model.C_Class).
			Select("id, name").
			Where("school_year_id=?", sy.Id).
			Order("id desc").Scan(&class).Error; err != nil {
			return nil, 0, errors.WithStack(err)
		}
		sy.Class = class
	}

	return schoolYears, count, nil
}

func (p *SchoolYear) Save(in *form.SaveSchoolYear) error {
	sy := model.SchoolYear{}
	tick := utee.Tick()
	if in.Id > 0 {
		if err := p.Ds.Db.First(&sy, in.Id).Error; err != nil {
			return errors.WithStack(err)
		}
	} else {
		sy.Ct = tick
	}
	if err := copier.Copy(&sy, in); err != nil {
		return errors.WithStack(err)
	}
	sy.Mt = tick

	err := p.Ds.Db.Save(&sy).Error
	if datasource.Duplicate(err) {
		return ut.NewValidateError("学年不可重复创建")
	}
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
