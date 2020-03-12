package hdl

import (
	"github.com/gin-gonic/gin"
	"github.com/quexer/utee"

	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/service"
	"liuyu/stu/pkg/web/hdl/form"
)

type Curriculum struct {
	Ds                *datasource.Ds
	CurriculumService *service.Curriculum
}

func (p *Curriculum) Mount(g *gin.RouterGroup) {
	g.GET("/", p.Get)
	g.GET("/list", p.List)
	g.GET("/listNameByIds", p.ListNameByIds)
	g.POST("/create", p.Create)
	g.POST("/update", p.Update)
	g.DELETE("/delete", p.Delete)
}

func (p *Curriculum) Get(c *gin.Context) {
	var req struct {
		Id int `form:"id"`
	}
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	var curriculum model.Curriculum
	utee.Chk(p.Ds.Db.First(&curriculum, req.Id).Error)

	c.JSON(200, curriculum)
}

func (p *Curriculum) List(c *gin.Context) {
	var req form.ListCurriculum
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	curriculums, total, err := p.CurriculumService.List(&req)
	utee.Chk(err)

	j := utee.J{
		"page":  req.Page,
		"limit": req.Limit,
		"total": total,
		"data":  curriculums,
	}
	c.JSON(200, j)
}

func (p *Curriculum) ListNameByIds(c *gin.Context) {
	var req struct {
		Ids []int `form:"ids"`
	}
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}

	if len(req.Ids) == 0 {
		c.JSON(200, utee.J{
			"data": []string{},
		})
		return
	}

	names, err := p.CurriculumService.ListNameByIds(req.Ids)
	utee.Chk(err)

	c.JSON(200, utee.J{
		"data": names,
	})
}

func (p *Curriculum) Create(c *gin.Context) {
	var req form.SaveCurriculum
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	val := func(req *form.SaveCurriculum) (bool, string) {
		if req.Name == "" {
			return false, "课程名称不能为空"
		}
		return true, ""
	}
	ok, str := val(&req)
	if !ok {
		c.String(400, str)
		return
	}
	utee.Chk(p.CurriculumService.Save(&req))
	c.String(200, "OK")
}

func (p *Curriculum) Update(c *gin.Context) {
	var req form.SaveCurriculum
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	val := func(req *form.SaveCurriculum) (bool, string) {
		if req.Id == 0 {
			return false, "Id不能为空"
		}
		if req.Name == "" {
			return false, "课程名称不能为空"
		}
		return true, ""
	}
	ok, str := val(&req)
	if !ok {
		c.String(400, str)
		return
	}
	utee.Chk(p.CurriculumService.Save(&req))
	c.String(200, "OK")
}

func (p *Curriculum) Delete(c *gin.Context) {
	var req struct {
		Id int `form:"id"`
	}
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}

	utee.Chk(p.Ds.Db.Where("id=?", req.Id).Delete(&model.Curriculum{}).Error)
	c.JSON(200, "ok")
}
