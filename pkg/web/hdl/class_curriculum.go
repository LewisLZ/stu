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

type ClassCurriculum struct {
	Ds                     *datasource.Ds
	ClassCurriculumService *service.ClassCurriculum
}

func (p *ClassCurriculum) Mount(g *gin.RouterGroup) {
	year := g.Group("/year")
	{
		year.GET("/list", p.YearList)
		year.POST("/create", p.YearSave(false))
		year.POST("/update", p.YearSave(false))
		year.DELETE("/delete", p.YearDelete)
	}

}

func (p *ClassCurriculum) YearList(c *gin.Context) {
	var req form.ListClassCurriculumYear
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	if req.ClassId == 0 {
		c.String(400, "课程不能为空")
		return
	}
	ccy, err := p.ClassCurriculumService.YearList(&req)
	if ut.IsValidateError(err) {
		c.String(400, err.Error())
		return
	}
	utee.Chk(err)

	c.JSON(200, utee.J{
		"data": ccy,
	})
}

func (p *ClassCurriculum) YearSave(update bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req form.SaveClassCurriculumYear
		if err := c.Bind(&req); err != nil {
			c.String(400, "参数错误")
			return
		}
		val := func(req *form.SaveClassCurriculumYear) (bool, string) {
			if update && req.Id == 0 {
				return false, "Id不能为空"
			}
			if req.ClassId == 0 {
				return false, "班级不能为空"
			}
			if req.Year == "" {
				return false, "年份不能为空"
			}
			if req.Pos < model.Pos_Up || req.Pos > model.Pos_Down {
				return false, "月份错误"
			}
			return true, ""
		}
		ok, str := val(&req)
		if !ok {
			c.String(400, str)
			return
		}
		utee.Chk(p.ClassCurriculumService.YearSave(&req))
		c.String(200, "OK")
	}
}

func (p *ClassCurriculum) YearDelete(c *gin.Context) {
	var req struct {
		Id int `form:"id"`
	}
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}

	utee.Chk(p.Ds.Db.Where("id=?", req.Id).Delete(&model.ClassCurriculumYear{}).Error)
	c.JSON(200, "ok")
}
