package hdl

import (
	"github.com/gin-gonic/gin"
	"github.com/quexer/utee"

	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/ut"
	"liuyu/stu/pkg/web/hdl/form"
)

type Pub struct {
	Ds *datasource.Ds
}

func (p *Pub) Mount(g *gin.RouterGroup) {
	g.POST("/login", p.Login)
}

// LoginHandler 登录
func (p *Pub) Login(c *gin.Context) {
	var req form.LoginRequest
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	if req.Account == "" || req.Passwd == "" {
		c.String(400, "account / passwd required")
		return
	}

	var user model.User

	if result := p.Ds.Db.Where("account = ? and passwd = ?", req.Account, ut.Passwd(req.Passwd)).First(&user); result.Error != nil {
		if result.RecordNotFound() {
			c.String(400, "登录失败")
			return
		}
		utee.Chk(result.Error)
	}

	j := gin.H{
		"id":     user.Id,
		"name":   user.Name,
		"mobile": user.Mobile,
	}

	utee.Chk(SetWebSession(c, &user))

	c.JSON(200, j)
}
