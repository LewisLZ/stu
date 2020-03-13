package hdl

import (
	"github.com/gin-gonic/gin"
	"github.com/quexer/utee"

	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/service"
	"liuyu/stu/pkg/web/hdl/form"
)

type Student struct {
	Ds             *datasource.Ds
	StudentService *service.Student
}

func (p *Student) Mount(g *gin.RouterGroup) {
	g.GET("", p.Get)
	g.GET("/list", p.List)
	g.POST("/create", p.Save(false))
	g.POST("/update", p.Save(true))
	g.DELETE("/delete", p.Delete)
}

func (p *Student) Get(c *gin.Context) {
	var req struct {
		Id int `form:"id"`
	}
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}

	stu, err := p.StudentService.Get(req.Id)
	utee.Chk(err)

	c.JSON(200, stu)
}

func (p *Student) List(c *gin.Context) {
	var req form.ListStudent
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}

	stus, total, err := p.StudentService.List(&req)
	utee.Chk(err)

	c.JSON(200, utee.J{
		"page":  req.Page,
		"limit": req.Limit,
		"total": total,
		"data":  stus,
	})
}

func (p *Student) Save(update bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req form.SaveStudent
		if err := c.Bind(&req); err != nil {
			c.String(400, "参数错误")
			return
		}

		val := func(req *form.SaveStudent) (bool, string) {
			if update && req.Id == 0 {
				return false, "Id不能为空"
			}
			if req.Name == "" {
				return false, "姓名不能为空"
			}
			if req.Mobile == "" {
				return false, "电话不能为空"
			}
			if req.Sex < model.SexBoy && req.Sex > model.SexGirl {
				return false, "性别错误"
			}
			if req.Birthday == "" {
				return false, "出生日期不能为空"
			}
			if req.IntakeTime == "" {
				return false, "入学时间不能为空"
			}
			if req.Address == "" {
				return false, "地址不能为空"
			}
			if req.ClassId == 0 {
				return false, "班级不能为空"
			}
			return true, ""
		}

		ok, str := val(&req)
		if !ok {
			c.String(400, str)
			return
		}

		utee.Chk(p.StudentService.Save(&req))

		c.String(200, "OK")
	}
}

func (p *Student) Delete(c *gin.Context) {
	var req struct {
		Id int `form:"id"`
	}
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}

	utee.Chk(p.Ds.Db.Where("id=?", req.Id).Delete(&model.Student{}).Error)
	c.JSON(200, "ok")
}
