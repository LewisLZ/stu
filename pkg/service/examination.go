package service

import (
	"time"

	"github.com/jinzhu/copier"
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
