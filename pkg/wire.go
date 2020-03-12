// +build wireinject

package pkg

import (
	"github.com/google/wire"

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
	datasource.CreateDs,
	webSet,
	helSet,
	srvSet,
)

var webSet = wire.NewSet(
	wire.Struct(new(web.Web), "*"),
)

var helSet = wire.NewSet(
	wire.Struct(new(hdl.Hdl), "*"),
	wire.Struct(new(hdl.Stu), "*"),
	wire.Struct(new(hdl.Pub), "*"),
	wire.Struct(new(hdl.Teacher), "*"),
	wire.Struct(new(hdl.User), "*"),
	wire.Struct(new(hdl.Mid), "*"),
	wire.Struct(new(hdl.Curriculum), "*"),
	wire.Struct(new(hdl.Class), "*"),
)

var srvSet = wire.NewSet(
	wire.Struct(new(service.Teacher), "*"),
	wire.Struct(new(service.Curriculum), "*"),
	wire.Struct(new(service.Class), "*"),
)
