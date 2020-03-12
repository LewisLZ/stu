package datasource

import (
	"github.com/quexer/utee"

	"liuyu/stu/pkg/model"
)

func initAndMigration(ds *Ds) {
	db := ds.Db
	err := db.AutoMigrate(
		model.Stu{},
		model.User{},
		model.Teacher{},
		model.Curriculum{},
		model.Class{},
		model.TeacherClass{},
		model.TeacherCurriculum{},
	).Error
	utee.Chk(err)
}
