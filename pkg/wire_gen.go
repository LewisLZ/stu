// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package pkg

import (
	"github.com/google/wire"
	"liuyu/stu/pkg/dao"
	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/service"
	"liuyu/stu/pkg/web"
	"liuyu/stu/pkg/web/hdl"
)

// Injectors from wire.go:

func New() (*App, error) {
	opt := defaultDsOpt()
	ds := datasource.CreateDs(opt)
	student := &service.Student{
		Ds: ds,
	}
	hdlStudent := &hdl.Student{
		Ds:             ds,
		StudentService: student,
	}
	class := &dao.Class{}
	curriculum := &dao.Curriculum{}
	teacher := &service.Teacher{
		Ds:            ds,
		ClassDao:      class,
		CurriculumDao: curriculum,
	}
	hdlTeacher := &hdl.Teacher{
		Ds:             ds,
		TeacherService: teacher,
	}
	pub := &hdl.Pub{
		Ds: ds,
	}
	user := &service.User{
		Ds: ds,
	}
	hdlUser := &hdl.User{
		Ds:          ds,
		UserService: user,
	}
	mid := &hdl.Mid{
		Ds: ds,
	}
	serviceCurriculum := &service.Curriculum{
		Ds: ds,
	}
	hdlCurriculum := &hdl.Curriculum{
		Ds:                ds,
		CurriculumService: serviceCurriculum,
	}
	serviceClass := &service.Class{
		Ds: ds,
	}
	hdlClass := &hdl.Class{
		Ds:           ds,
		ClassService: serviceClass,
	}
	schoolYear := &service.SchoolYear{
		Ds: ds,
	}
	hdlSchoolYear := &hdl.SchoolYear{
		Ds:                ds,
		SchoolYearService: schoolYear,
	}
	classCurriculum := &dao.ClassCurriculum{}
	serviceClassCurriculum := &service.ClassCurriculum{
		Ds:                 ds,
		ClassCurriculumDao: classCurriculum,
	}
	hdlClassCurriculum := &hdl.ClassCurriculum{
		Ds:                     ds,
		ClassCurriculumService: serviceClassCurriculum,
	}
	examination := &service.Examination{
		Ds: ds,
	}
	hdlExamination := &hdl.Examination{
		Ds:                 ds,
		ExaminationService: examination,
	}
	achievement := &service.Achievement{
		Ds: ds,
	}
	hdlAchievement := &hdl.Achievement{
		AchievementService: achievement,
	}
	hdlHdl := &hdl.Hdl{
		Ds:              ds,
		Student:         hdlStudent,
		Teacher:         hdlTeacher,
		Pub:             pub,
		User:            hdlUser,
		Mid:             mid,
		Curriculum:      hdlCurriculum,
		Class:           hdlClass,
		SchoolYear:      hdlSchoolYear,
		ClassCurriculum: hdlClassCurriculum,
		Examination:     hdlExamination,
		Achievement:     hdlAchievement,
	}
	webWeb := &web.Web{
		Hdl: hdlHdl,
		Ds:  ds,
	}
	webOpt := defaultWebOpt()
	app := &App{
		Web:    webWeb,
		WebOpt: webOpt,
	}
	return app, nil
}

// wire.go:

var appSet = wire.NewSet(wire.Struct(new(App), "*"), defaultDsOpt,
	defaultWebOpt, datasource.CreateDs, webSet,
	helSet,
	srvSet,
	daoSet,
)

var webSet = wire.NewSet(wire.Struct(new(web.Web), "*"))

var helSet = wire.NewSet(wire.Struct(new(hdl.Hdl), "*"), wire.Struct(new(hdl.Student), "*"), wire.Struct(new(hdl.Pub), "*"), wire.Struct(new(hdl.Teacher), "*"), wire.Struct(new(hdl.User), "*"), wire.Struct(new(hdl.Mid), "*"), wire.Struct(new(hdl.Curriculum), "*"), wire.Struct(new(hdl.Class), "*"), wire.Struct(new(hdl.SchoolYear), "*"), wire.Struct(new(hdl.ClassCurriculum), "*"), wire.Struct(new(hdl.Examination), "*"), wire.Struct(new(hdl.Achievement), "*"))

var srvSet = wire.NewSet(wire.Struct(new(service.Teacher), "*"), wire.Struct(new(service.Curriculum), "*"), wire.Struct(new(service.Class), "*"), wire.Struct(new(service.Student), "*"), wire.Struct(new(service.SchoolYear), "*"), wire.Struct(new(service.ClassCurriculum), "*"), wire.Struct(new(service.Examination), "*"), wire.Struct(new(service.Achievement), "*"), wire.Struct(new(service.User), "*"))

var daoSet = wire.NewSet(wire.Struct(new(dao.Class), "*"), wire.Struct(new(dao.Curriculum), "*"), wire.Struct(new(dao.ClassCurriculum), "*"))
