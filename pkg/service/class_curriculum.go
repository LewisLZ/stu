package service

import (
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/quexer/utee"

	"liuyu/stu/pkg/dao"
	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/ut"
	"liuyu/stu/pkg/web/hdl/form"
)

type ClassCurriculum struct {
	Ds                 *datasource.Ds
	ClassCurriculumDao *dao.ClassCurriculum
}

func (p *ClassCurriculum) YearList(req *form.ListClassCurriculumYear) ([]*model.ClassCurriculumYear, error) {
	var ccy []*model.ClassCurriculumYear
	if err := p.Ds.Db.Table("class_curriculum_year cc").
		Select("cc.id, cc.year, cc.pos, c.name class_name, cc.class_id").
		Joins("LEFT JOIN class c ON cc.class_id=c.id").
		Where("cc.class_id=?", req.ClassId).
		Order("cc.id desc").
		Scan(&ccy).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	for _, v := range ccy {
		v.ClassName = ut.MakeClassName(v.ClassName, v.Year, v.Pos)
		cc, err := p.ClassCurriculumDao.ListClassCurriculumWithYearId(p.Ds.Db, v.Id)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		v.ClassCurriculum = cc

		tmp, err := time.ParseInLocation("2006", v.Year, time.Local)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		v.YearTmp = utee.Tick(tmp)
	}

	return ccy, nil
}

func (p *ClassCurriculum) YearSave(req *form.SaveClassCurriculumYear) error {

	var school struct {
		Year string
		Pos  model.Pos
	}
	if err := p.Ds.Db.Table("class c, school_year s").
		Select("s.year, s.pos").
		Where("c.id=?", req.ClassId).
		Scan(&school).Error; err != nil {
		return errors.WithStack(err)
	}
	reqTmp, err := time.ParseInLocation("2006", req.Year, time.Local)
	if err != nil {
		return errors.WithStack(err)
	}
	schoolTmp, err := time.ParseInLocation("2006", req.Year, time.Local)
	if err != nil {
		return errors.WithStack(err)
	}

	if reqTmp.Before(schoolTmp) {
		return ut.NewValidateError("年份不能小于班级")
	}
	if reqTmp.Equal(schoolTmp) && req.Pos < school.Pos {
		return ut.NewValidateError("年份等于班级时月份不能小于班级")
	}

	var cc model.ClassCurriculumYear
	pick := utee.Tick()
	if req.Id != 0 {
		if err := p.Ds.Db.First(&cc, req.Id).Error; err != nil {
			return errors.WithStack(err)
		}
	} else {
		cc.Ct = pick
	}
	if err := copier.Copy(&cc, req); err != nil {
		return errors.WithStack(err)
	}
	cc.Mt = pick

	if err := p.Ds.Db.Save(&cc).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
