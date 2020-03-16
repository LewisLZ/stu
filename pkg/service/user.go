package service

import (
	"github.com/quexer/utee"

	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/ut"
	"liuyu/stu/pkg/web/hdl/form"
)

type User struct {
	Ds *datasource.Ds
}

func (p *User) List(in *form.ListUser) ([]*model.User, int, error) {
	db := p.Ds.Db.Model(&model.User{}).Select("id, mobile, name, type")
	if in.Mobile != "" {
		db = db.Where("mobile like ?", "%"+in.Mobile+"%")
	}
	if in.Name != "" {
		db = db.Where("name like ?", "%"+in.Name+"%")
	}

	var count int
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var users []*model.User
	_, limit, offset := ut.MakePager(in.Page, in.Limit, 10)
	if err := db.Limit(limit).Offset(offset).Order("id desc").Scan(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (p *User) Save(in *form.SaveUser) error {
	var user model.User
	tick := utee.Tick()
	if in.Id != 0 {
		if err := p.Ds.Db.First(&user, in.Id).Error; err != nil {
			return err
		}
	} else {
		user.Ct = tick
	}
	user.Mt = tick
	user.Mobile = in.Mobile
	user.Name = in.Name
	if in.Passwd != "" {
		user.Passwd = ut.Passwd(in.Passwd)
	}

	if err := p.Ds.Db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}
