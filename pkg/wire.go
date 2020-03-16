// +build wireinject

package pkg

import (
	"github.com/google/wire"

	"liuyu/stu/pkg/dao"
	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/service"
	"liuyu/stu/pkg/web"
	"liuyu/stu/pkg/web/hdl"
)

func New() (*App, error) {
	panic(wire.Build(appSet))
}

var appSet = wire.NewSet(
	wire.Struct(new(App), "*"),
	defaultDsOpt,
	defaultWebOpt,
	datasource.CreateDs,
	webSet,
	helSet,
	srvSet,
	daoSet,
)

var webSet = wire.NewSet(
	wire.Struct(new(web.Web), "*"),
)

var helSet = wire.NewSet(
	wire.Struct(new(hdl.Hdl), "*"),
	wire.Struct(new(hdl.Student), "*"),
	wire.Struct(new(hdl.Pub), "*"),
	wire.Struct(new(hdl.Teacher), "*"),
	wire.Struct(new(hdl.User), "*"),
	wire.Struct(new(hdl.Mid), "*"),
	wire.Struct(new(hdl.Curriculum), "*"),
	wire.Struct(new(hdl.Class), "*"),
	wire.Struct(new(hdl.SchoolYear), "*"),
	wire.Struct(new(hdl.ClassCurriculum), "*"),
	wire.Struct(new(hdl.Examination), "*"),
	wire.Struct(new(hdl.Achievement), "*"),
)

var srvSet = wire.NewSet(
	wire.Struct(new(service.Teacher), "*"),
	wire.Struct(new(service.Curriculum), "*"),
	wire.Struct(new(service.Class), "*"),
	wire.Struct(new(service.Student), "*"),
	wire.Struct(new(service.SchoolYear), "*"),
	wire.Struct(new(service.ClassCurriculum), "*"),
	wire.Struct(new(service.Examination), "*"),
	wire.Struct(new(service.Achievement), "*"),
	wire.Struct(new(service.User), "*"),
)

var daoSet = wire.NewSet(
	wire.Struct(new(dao.Class), "*"),
	wire.Struct(new(dao.Curriculum), "*"),
	wire.Struct(new(dao.ClassCurriculum), "*"),
)
