package hdl

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/quexer/utee"

	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/service"
	"liuyu/stu/pkg/web/hdl/form"
)

type Teacher struct {
	Ds             *datasource.Ds
	TeacherService *service.Teacher
}

func (p *Teacher) Mount(g *gin.RouterGroup) {
	g.GET("/", p.Get)
	g.GET("/list", p.List)
	g.POST("/create", p.Create)
	g.POST("/update", p.Update)
	g.DELETE("/delete", p.Delete)
}

func (p *Teacher) Get(c *gin.Context) {
	var req struct {
		Id int `form:"id"`
	}
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}

	teacher, err := p.TeacherService.Get(req.Id)
	utee.Chk(err)

	c.JSON(200, teacher)
}

func (p *Teacher) List(c *gin.Context) {
	var req form.ListTeacher
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	if req.Sex < model.SexUnknown || req.Sex > model.SexGirl {
		c.String(400, "性别错误")
		return
	}

	teachers, count, err := p.TeacherService.List(&req)
	utee.Chk(err)

	j := utee.J{
		"data":  teachers,
		"total": count,
		"page":  req.Page,
		"limit": req.Limit,
	}
	c.JSON(200, j)
}

func (p *Teacher) Create(c *gin.Context) {
	var req form.SaveTeacher
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	val := func(teacher *form.SaveTeacher) (bool, string) {
		if teacher.Name == "" {
			return false, "姓名不能为空"
		}
		if teacher.Mobile == "" {
			return false, "电话不能为空"
		}
		if teacher.Sex != 1 && teacher.Sex != 2 {
			return false, "性别错误"
		}
		if len(teacher.CurriculumIds) == 0 {
			return false, "课程不能为空"
		}
		if len(teacher.ClassIds) == 0 {
			return false, "班级不能为空"
		}
		return true, ""
	}
	ok, s := val(&req)
	if !ok {
		c.String(400, s)
		return
	}
	utee.Chk(p.TeacherService.Save(&req))

	c.String(200, "OK")
}

func (p *Teacher) Update(c *gin.Context) {
	var req form.SaveTeacher
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	val := func(teacher *form.SaveTeacher) (bool, string) {
		if teacher.Id == 0 {
			return false, "Id不能为空"
		}
		if teacher.Name == "" {
			return false, "姓名不能为空"
		}
		if teacher.Mobile == "" {
			return false, "电话不能为空"
		}
		if teacher.Sex != 1 && teacher.Sex != 2 {
			return false, "性别错误"
		}
		if len(teacher.CurriculumIds) == 0 {
			return false, "课程不能为空"
		}
		if len(teacher.ClassIds) == 0 {
			return false, "班级不能为空"
		}
		return true, ""
	}
	ok, s := val(&req)
	if !ok {
		c.String(400, s)
		return
	}
	utee.Chk(p.TeacherService.Save(&req))

	c.String(200, "OK")
}

func (p *Teacher) Delete(c *gin.Context) {
	var req struct {
		Id int `form:"id"`
	}
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}

	err := datasource.RunTransaction(p.Ds.Db, func(tx *gorm.DB) error {
		if err := p.Ds.Db.Where("id=?", req.Id).Delete(&model.Teacher{}).Error; err != nil {
			return err
		}
		if err := tx.Where("teacher_id=?", req.Id).Delete(&model.TeacherClass{}).Error; err != nil {
			return err
		}

		if err := tx.Where("teacher_id=?", req.Id).Delete(&model.TeacherCurriculum{}).Error; err != nil {
			return err
		}
		return nil
	})
	utee.Chk(err)
	c.JSON(200, "ok")
}
