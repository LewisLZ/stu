package service

import (
	"github.com/jinzhu/copier"
	"github.com/quexer/utee"

	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/web/hdl/form"
)

type Class struct {
	Ds *datasource.Ds
}

func (p *Class) List(parentId int) ([]*model.Class, error) {
	var clasies []*model.Class

	err := p.Ds.Db.Model(&model.Class{}).
		Select("id, parent_id, name").
		Where("parent_id=?", parentId).
		Scan(&clasies).Error

	if err != nil {
		return nil, err
	}

	for _, class := range clasies {
		cls, err := p.List(class.Id)
		if err != nil {
			return nil, err
		}
		class.Children = cls

		teachers := []*model.Teacher{}
		if err := p.Ds.Db.Table(model.C_Teacher+" t").
			Joins("LEFT JOIN teacher_class tc ON tc.teacher_id=t.id").
			Select("t.id, t.name, t.mobile").
			Where("tc.class_id=?", class.Id).Order("t.id desc").Scan(&teachers).Error; err != nil {
			return nil, err
		}
		class.Teacher = teachers
	}

	return clasies, nil
}

func (p *Class) ListNameByIds(ids []int) ([]string, error) {
	var names []string

	if err := p.Ds.Db.Raw(`select concat(c2.name, '-', c1.name) names
			from class c1, class c2 
			where c1.id in (?) and c1.parent_id=c2.id`, ids).Order("c1.id desc").Pluck("names", &names).
		Error; err != nil {
		return nil, err
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
