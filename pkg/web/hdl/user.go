package hdl

import (
	"github.com/gin-gonic/gin"
	"github.com/quexer/utee"

	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/service"
	"liuyu/stu/pkg/web/hdl/form"
)

type User struct {
	Ds          *datasource.Ds
	UserService *service.User
}

func (p *User) Mount(g *gin.RouterGroup) {
	g.GET("/uc", p.Uc)
	g.GET("/list", p.List)
	g.POST("/create", p.Save(false))
	g.POST("/update", p.Save(true))
	g.DELETE("/delete", p.Delete)
}

func (p *User) Uc(c *gin.Context) {
	user := MustGetUser(c)
	j := gin.H{
		"id":     user.Id,
		"name":   user.Name,
		"mobile": user.Mobile,
		"type":   user.Type,
	}
	c.JSON(200, j)
}

func (p *User) List(c *gin.Context) {
	var req form.ListUser
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	users, total, err := p.UserService.List(&req)
	utee.Chk(err)

	c.JSON(200, utee.J{
		"page":  req.Page,
		"limit": req.Limit,
		"total": total,
		"data":  users,
	})
}

func (p *User) Save(update bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req form.SaveUser
		if err := c.Bind(&req); err != nil {
			c.String(400, "参数错误")
			return
		}
		val := func(req *form.SaveUser) (bool, string) {
			if update && req.Id == 0 {
				return false, "Id不能为空"
			}
			if req.Name == "" {
				return false, "姓名不能为空"
			}
			if req.Mobile == "" {
				return false, "电话不能为空"
			}
			if !update && req.Passwd == "" {
				return false, "电话不能为空"
			}
			return true, ""
		}

		ok, str := val(&req)
		if !ok {
			c.String(400, str)
			return
		}

		utee.Chk(p.UserService.Save(&req))

		c.String(200, "OK")
	}
}

func (p *User) Delete(c *gin.Context) {
	var req struct {
		Id int `form:"id"`
	}
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}

	var user model.User
	err := p.Ds.Db.First(&user, req.Id).Error
	if datasource.NotFound(err) {
		c.String(400, "没找到管理员")
		return
	}
	utee.Chk(err)
	if user.Type == model.UserType_Super {
		c.String(400, "不能删除管理员")
		return
	}

	utee.Chk(p.Ds.Db.Where("id=?", req.Id).Delete(&model.User{}).Error)

	c.String(200, "OK")
}
