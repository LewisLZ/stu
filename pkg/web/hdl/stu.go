package hdl

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Stu struct {
}

func (p *Stu) Mount(g *gin.RouterGroup) {
	g.GET("/list", p.List)
}

func (p *Stu) List(c *gin.Context) {
	fmt.Println("========")

	c.String(200, "oooooooo")
}
