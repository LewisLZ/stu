package hdl

import (
	"github.com/gin-gonic/gin"
	"github.com/quexer/utee"

	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/service"
	"liuyu/stu/pkg/ut"
	"liuyu/stu/pkg/web/hdl/form"
)

type Examination struct {
	Ds                 *datasource.Ds
	ExaminationService *service.Examination
}

func (p *Examination) Mount(g *gin.RouterGroup) {
	g.GET("/list", p.List)
	g.POST("/create", p.Save(false))
	g.POST("/update", p.Save(true))
	g.DELETE("/delete", p.Delete)

	cla := g.Group("/class")
	{
		cla.GET("/list", p.ClassList)
		cla.POST("/create", p.ClassCreate)
		cla.DELETE("/delete", p.ClassDelete)
	}
	cc := g.Group("/classcurriculum")
	{
		cc.POST("/create", p.ClassCurriculumCreate)
	}
}

func (p *Examination) List(c *gin.Context) {
	var req form.ListExamination
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	if (req.StartTime != "" && req.EndTime == "") || (req.StartTime == "" && req.EndTime != "") {
		c.String(400, "开始和结束时间必须同时有值")
		return
	}
	es, total, err := p.ExaminationService.List(&req)
	if ut.IsValidateError(err) {
		c.String(400, err.Error())
		return
	}

	c.JSON(200, utee.J{
		"page":  req.Page,
		"limit": req.Limit,
		"total": total,
		"data":  es,
	})
}

func (p *Examination) Save(update bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req form.SaveExamination
		if err := c.Bind(&req); err != nil {
			c.String(400, "参数错误")
			return
		}
		val := func(req *form.SaveExamination) (bool, string) {
			if update && req.Id == 0 {
				return false, "Id不能为空"
			}
			if req.Name == "" {
				return false, "名称不能为空"
			}
			if req.StartTime == "" {
				return false, "开始时间不能为空"
			}
			return true, ""
		}
		ok, str := val(&req)
		if !ok {
			c.String(400, str)
			return
		}
		err := p.ExaminationService.Save(&req)
		if ut.IsValidateError(err) {
			c.String(400, err.Error())
			return
		}
		utee.Chk(err)
		c.String(200, "OK")
	}
}

func (p *Examination) Delete(c *gin.Context) {
	var req struct {
		Id int `form:"id"`
	}
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}

	utee.Chk(p.Ds.Db.Where("id=?", req.Id).Delete(&model.Examination{}).Error)
	c.JSON(200, "ok")
}

func (p *Examination) ClassList(c *gin.Context) {
	var req form.ListExaminationClass
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	if req.ExaminationId == 0 {
		c.String(400, "考试Id不能为空")
		return
	}
	ecs, canEdit, err := p.ExaminationService.ClassList(&req)
	utee.Chk(err)

	c.JSON(200, utee.J{
		"data":     ecs,
		"can_edit": canEdit,
	})
}

func (p *Examination) ClassCreate(c *gin.Context) {
	var req form.SaveExaminationClass
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	val := func(req *form.SaveExaminationClass) (bool, string) {
		if req.ExaminationId == 0 {
			return false, "考试Id不能为空"
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
	err := p.ExaminationService.ClassSave(&req)
	if ut.IsValidateError(err) {
		c.String(400, err.Error())
		return
	}
	utee.Chk(err)
	c.String(200, "OK")
}

func (p *Examination) ClassDelete(c *gin.Context) {
	var req struct {
		Id int `form:"id"`
	}
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}

	utee.Chk(p.Ds.Db.Where("id=?", req.Id).Delete(&model.ExaminationClass{}).Error)
	c.JSON(200, "ok")
}

func (p *Examination) ClassCurriculumCreate(c *gin.Context) {
	var req form.SaveExaminationClassCurriculum
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	val := func(req *form.SaveExaminationClassCurriculum) (bool, string) {
		if req.ExaminationClassId == 0 {
			return false, "考试Id不能为空"
		}
		return true, ""
	}
	ok, str := val(&req)
	if !ok {
		c.String(400, str)
		return
	}
	err := p.ExaminationService.ClassCurriculumSave(&req)
	if ut.IsValidateError(err) {
		c.String(400, err.Error())
		return
	}
	utee.Chk(err)
	c.String(200, "OK")
}
