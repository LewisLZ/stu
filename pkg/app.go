package pkg

import (
	"liuyu/stu/pkg/web"
)

type App struct {
	Web *web.Web
}

func (p *App) Run() {
	p.Web.CreateWebServer(":3000", "./public")
}
