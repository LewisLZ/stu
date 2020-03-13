package hdl

import (
	"github.com/gin-gonic/gin"
	"github.com/quexer/utee"

	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/service"
	"liuyu/stu/pkg/web/hdl/form"
)

type SchoolYear struct {
	Ds                *datasource.Ds
	SchoolYearService *service.SchoolYear
}

func (p *SchoolYear) Mount(g *gin.RouterGroup) {
	g.GET("/list", p.List)
	g.GET("/listNameByIds", p.ListNameByIds)
	g.POST("/create", p.Save(false))
	g.POST("/update", p.Save(true))
	g.DELETE("/delete", p.Delete)
}

func (p *SchoolYear) List(c *gin.Context) {
	var req form.ListSchoolYear
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	sy, total, err := p.SchoolYearService.List(&req)
	utee.Chk(err)

	c.JSON(200, utee.J{
		"page":  req.Page,
		"limit": req.Limit,
		"total": total,
		"data":  sy,
	})
}

func (p *SchoolYear) ListNameByIds(c *gin.Context) {
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

	names, err := p.SchoolYearService.ListNameByIds(req.Ids)
	utee.Chk(err)

	c.JSON(200, utee.J{
		"data": names,
	})
}

func (p *SchoolYear) Save(update bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req form.SaveSchoolYear
		if err := c.Bind(&req); err != nil {
			c.String(400, "参数错误")
			return
		}
		val := func(req *form.SaveSchoolYear) (bool, string) {
			if update && req.Id == 0 {
				return false, "Id不能为空"
			}
			if req.Year == "" {
				return false, "年不能为空"
			}
			if req.Pos < model.Pos_Up || req.Pos > model.Pos_Down {
				return false, "入学月份错误"
			}
			return true, ""
		}

		ok, str := val(&req)
		if !ok {
			c.String(400, str)
			return
		}

		utee.Chk(p.SchoolYearService.Save(&req))

		c.String(200, "OK")
	}
}

func (p *SchoolYear) Delete(c *gin.Context) {
	var req struct {
		Id int `form:"id"`
	}
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}

	var count int
	utee.Chk(p.Ds.Db.Model(model.Class{}).Where("school_year_id=?", req.Id).Count(&count).Error)

	if count > 0 {
		c.String(400, "此学年班级不可删除")
		return
	}

	utee.Chk(p.Ds.Db.Where("id=?", req.Id).Delete(&model.SchoolYear{}).Error)
	c.JSON(200, "ok")
}
