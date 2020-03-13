package hdl

import (
	"github.com/gin-gonic/gin"
	"github.com/quexer/utee"

	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/service"
	"liuyu/stu/pkg/web/hdl/form"
)

type Class struct {
	Ds           *datasource.Ds
	ClassService *service.Class
}

func (p *Class) Mount(g *gin.RouterGroup) {
	g.GET("", p.Get)
	g.GET("/list", p.List)
	g.GET("/listNameByIds", p.ListNameByIds)
	g.POST("/create", p.Save(false))
	g.POST("/update", p.Save(true))
	g.DELETE("/delete", p.Delete)
}

func (p *Class) Get(c *gin.Context) {
	var req struct {
		Id int `form:"id"`
	}
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	res, err := p.ClassService.Get(req.Id)
	utee.Chk(err)

	c.JSON(200, res)
}
func (p *Class) List(c *gin.Context) {
	var req form.ListClass
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	res, total, err := p.ClassService.List(&req)
	utee.Chk(err)

	c.JSON(200, utee.J{
		"page":  req.Page,
		"limit": req.Limit,
		"total": total,
		"data":  res,
	})
}

func (p *Class) ListNameByIds(c *gin.Context) {
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

	names, err := p.ClassService.ListNameByIds(req.Ids)
	utee.Chk(err)

	c.JSON(200, utee.J{
		"data": names,
	})
}

func (p *Class) Save(update bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req form.SaveClass
		if err := c.Bind(&req); err != nil {
			c.String(400, "参数错误")
			return
		}
		val := func(req *form.SaveClass) (bool, string) {
			if update && req.Id == 0 {
				return false, "Id不能为空"
			}
			if req.Name == "" {
				return false, "名称不能为空"
			}
			if req.SchoolYearId == 0 {
				return false, "学年不能为空"
			}
			return true, ""
		}

		ok, str := val(&req)
		if !ok {
			c.String(400, str)
			return
		}

		utee.Chk(p.ClassService.Save(&req))

		c.String(200, "OK")
	}
}

func (p *Class) Delete(c *gin.Context) {
	var req struct {
		Id int `form:"id"`
	}
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}

	var count int
	utee.Chk(p.Ds.Db.Model(model.Class{}).Where("parent_id=?", req.Id).Count(&count).Error)

	if count > 0 {
		c.String(400, "有下级班级不可删除")
		return
	}

	utee.Chk(p.Ds.Db.Where("id=?", req.Id).Delete(&model.Class{}).Error)
	c.JSON(200, "ok")
}
