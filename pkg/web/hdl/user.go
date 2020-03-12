package hdl

import (
	"github.com/gin-gonic/gin"

	"liuyu/stu/pkg/datasource"
)

type User struct {
	Ds *datasource.Ds
}

func (p *User) Mount(g *gin.RouterGroup) {
	g.GET("/uc", p.Uc)
}

func (p *User) Uc(c *gin.Context) {
	user := MustGetUser(c)
	j := gin.H{
		"id":     user.Id,
		"name":   user.Name,
		"mobile": user.Mobile,
	}
	c.JSON(200, j)
}
