package service

import (
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"github.com/quexer/utee"

	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/ut"
	"liuyu/stu/pkg/web/hdl/form"
)

type Examination struct {
	Ds *datasource.Ds
}

func (p *Examination) List(in *form.ListExamination) ([]*model.Examination, int, error) {
	db := p.Ds.Db.Model(&model.Examination{}).Select("id, name, start_time, remark")
	if in.Name != "" {
		db = db.Where("name like ?", "%"+in.Name+"%")
	}
	if in.StartTime != "" && in.EndTime != "" {
		start, err := time.ParseInLocation("2006-01-02", in.StartTime, time.Local)
		if err != nil {
			return nil, 0, ut.NewValidateError("开始时间格式不对")
		}
		end, err := time.ParseInLocation("2006-01-02", in.EndTime, time.Local)
		if err != nil {
			return nil, 0, ut.NewValidateError("结束时间格式不对")
		}
		db = db.Where("start_time>=? AND start_time<=?", utee.Tick(start), utee.Tick(end))
	}

	var count int
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	var examinations []*model.Examination
	_, limit, offset := ut.MakePager(in.Page, in.Limit, 10)
	if err := db.Limit(limit).Offset(offset).Order("id desc").Scan(&examinations).Error; err != nil {
		return nil, 0, err
	}

	for _, e := range examinations {
		var count int
		if err := p.Ds.Db.Model(&model.ExaminationClass{}).
			Where("examination_id=?", e.Id).
			Count(&count).Error; err != nil {
			return nil, 0, err
		}
		e.ExaminationItemCount = count
	}

	return examinations, count, nil
}

func (p *Examination) Save(in *form.SaveExamination) error {
	var examination model.Examination
	pick := utee.Tick()
	if in.Id != 0 {
		if err := p.Ds.Db.First(&examination, in.Id).Error; err != nil {
			return err
		}
	} else {
		examination.Ct = pick
	}
	if err := copier.Copy(&examination, in); err != nil {
		return err
	}
	start, err := time.ParseInLocation("2006-01-02", in.StartTime, time.Local)
	if err != nil {
		return ut.NewValidateError("开始时间格式不对")
	}
	examination.Mt = pick
	examination.StartTime = utee.Tick(start)

	if err := p.Ds.Db.Save(&examination).Error; err != nil {
		return err
	}
	return nil
}

func (p *Examination) ClassList(in *form.ListExaminationClass) ([]*model.ExaminationClass, error) {
	var ecs []*model.ExaminationClass

	if err := p.Ds.Db.Table("examination_class ec").
		Select("ec.id, ec.examination_id, ec.class_id, c.name class_name, s.year, s.pos").
		Joins("LEFT JOIN class c ON c.id=ec.class_id").
		Joins("LEFT JOIN school_year s ON s.id=c.school_year_id").
		Where("ec.examination_id=?", in.ExaminationId).
		Order("ec.id desc").Scan(&ecs).Error; err != nil {
		return nil, err
	}

	for _, ec := range ecs {
		ec.ClassName = ut.MakeClassName(ec.ClassName, ec.Year, ec.Pos)

		var ecc []*model.ExaminationClassCurriculum

		if err := p.Ds.Db.Table("examination_class_curriculum ecc").
			Joins("LEFT JOIN class_curriculum cc ON cc.id=ecc.class_curriculum_id").
			Joins("LEFT JOIN class_curriculum_year ccy ON ccy.id=cc.cc_year_id").
			Joins("LEFT JOIN curriculum c ON c.id=cc.curriculum_id").
			Select("ecc.id, ecc.examination_class_id, ecc.class_curriculum_id, ccy.year class_curriculum_year, ccy.pos class_curriculum_pos, c.name class_curriculum_name").
			Where("ecc.examination_class_id=?", ec.Id).Scan(&ecc).Error; err != nil {
			return nil, err
		}

		for _, v := range ecc {
			v.ClassCurriculumName = ut.MakeClassName(v.ClassCurriculumName, v.ClassCurriculumYear, v.ClassCurriculumPos)
		}

		ec.ExaminationClassCurriculum = ecc
	}

	return ecs, nil
}

func (p *Examination) ClassSave(in *form.SaveExaminationClass) error {
	tick := utee.Tick()
	ec := model.ExaminationClass{}
	ec.ExaminationId = in.ExaminationId
	ec.ClassId = in.ClassId
	ec.Ct = tick
	ec.Mt = tick

	if err := p.Ds.Db.Create(&ec).Error; err != nil {
		return err
	}
	return nil
}

func (p *Examination) ClassCurriculumSave(in *form.SaveExaminationClassCurriculum) error {
	tick := utee.Tick()

	err := datasource.RunTransaction(p.Ds.Db, func(tx *gorm.DB) error {
		var ids []int
		if err := tx.Model(&model.ExaminationClassCurriculum{}).
			Select("class_curriculum_id").
			Where("examination_class_id=?", in.ExaminationClassId).
			Pluck("class_curriculum_id", &ids).Error; err != nil {
			return err
		}

		dataSet := mapset.NewSet()
		for _, item := range ids {
			dataSet.Add(item)
		}
		inSet := mapset.NewSet()
		for _, item := range in.ClassCurriculumIds {
			inSet.Add(item)
		}

		deleteIds := dataSet.Difference(inSet)
		addIds := inSet.Difference(dataSet)

		for _, id := range deleteIds.ToSlice() {
			if err := tx.Where("examination_class_id=? AND class_curriculum_id=?", in.ExaminationClassId, id).
				Delete(&model.ExaminationClassCurriculum{}).Error; err != nil {
				return err
			}
		}

		for _, id := range addIds.ToSlice() {
			ecc := model.ExaminationClassCurriculum{}
			ecc.ExaminationClassId = in.ExaminationClassId
			ecc.ClassCurriculumId = id.(int)
			ecc.Ct = tick
			ecc.Mt = tick
			if err := tx.Create(&ecc).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
