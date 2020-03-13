package datasource

import (
	"github.com/quexer/utee"

	"liuyu/stu/pkg/model"
)

func initAndMigration(ds *Ds) {
	db := ds.Db
	err := db.AutoMigrate(
		model.Student{},
		model.StudentCode{},
		model.User{},
		model.Teacher{},
		model.Curriculum{},
		model.SchoolYear{},
		model.Class{},
		model.TeacherClass{},
		model.TeacherCurriculum{},
		model.ClassCurriculumYear{},
		model.ClassCurriculum{},
	).Error
	utee.Chk(err)
}
