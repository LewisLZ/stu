package pkg

import (
	"liuyu/stu/pkg/ut"
	"liuyu/stu/pkg/web"
)

type App struct {
	Web    *web.Web
	WebOpt *ut.WebOpt
}

func (p *App) Run() {
	p.Web.CreateWebServer(p.WebOpt)
}
