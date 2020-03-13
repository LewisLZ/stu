package hdl

import (
	"github.com/gin-gonic/gin"
)

type Student struct {
}

func (p *Student) Mount(g *gin.RouterGroup) {
	g.GET("", p.Get)
	g.GET("/list", p.List)
	g.POST("/create", p.Create)
	g.POST("/update", p.Update)
	g.DELETE("/delete", p.Delete)
}

func (p *Student) Get(c *gin.Context) {
	c.String(200, "oooooooo")
}

func (p *Student) List(c *gin.Context) {
	c.String(200, "oooooooo")
}

func (p *Student) Create(c *gin.Context) {
	c.String(200, "oooooooo")
}

func (p *Student) Update(c *gin.Context) {
	c.String(200, "oooooooo")
}

func (p *Student) Delete(c *gin.Context) {
	c.String(200, "oooooooo")
}
