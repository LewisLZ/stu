package hdl

import (
	"github.com/gin-gonic/gin"
	"github.com/quexer/utee"

	"liuyu/stu/pkg/service"
	"liuyu/stu/pkg/ut"
	"liuyu/stu/pkg/web/hdl/form"
)

type Achievement struct {
	AchievementService *service.Achievement
}

func (p *Achievement) Mount(g *gin.RouterGroup) {
	g.GET("/list", p.List)
	g.GET("/listsocre", p.ListScore)
	g.GET("/listsearchsocre", p.ListSearchScore)
	g.GET("/liststudentscore", p.ListStudentScore)
	g.GET("/listcurriculumscore", p.ListCurriculumScore)
	g.GET("/listclassscore", p.ListClassScore)
	g.POST("/savesocre", p.SaveScore)
	g.POST("/archived", p.Archived)
}

func (p *Achievement) List(c *gin.Context) {
	var req form.ListAchievement
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	if (req.StartTime != "" && req.EndTime == "") || (req.StartTime == "" && req.EndTime != "") {
		c.String(400, "开始和结束时间必须同时有值")
		return
	}
	as, total, err := p.AchievementService.List(&req)
	if ut.IsValidateError(err) {
		c.String(400, err.Error())
		return
	}
	utee.Chk(err)

	c.JSON(200, utee.J{
		"page":  req.Page,
		"limit": req.Limit,
		"total": total,
		"data":  as,
	})
}

func (p *Achievement) ListScore(c *gin.Context) {
	var req form.ListAchievementSource
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	if req.ExaminationClassId == 0 {
		c.String(400, "考试班级Id不能为空")
		return
	}
	as, err := p.AchievementService.ListSource(&req)
	if ut.IsValidateError(err) {
		c.String(400, err.Error())
		return
	}
	utee.Chk(err)

	c.JSON(200, utee.J{
		"data": as,
	})
}

func (p *Achievement) ListStudentScore(c *gin.Context) {
	var req form.ListAchievementStudentScore
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	if req.StudentId == 0 {
		c.String(400, "学生Id不能为空")
		return
	}
	if req.ExaminationClassId == 0 {
		c.String(400, "考试Id不能为空")
		return
	}

	as, err := p.AchievementService.ListStudentScore(&req)
	if ut.IsValidateError(err) {
		c.String(400, err.Error())
		return
	}
	utee.Chk(err)

	c.JSON(200, utee.J{
		"data": as,
	})
}

func (p *Achievement) ListCurriculumScore(c *gin.Context) {
	var req form.ListAchievementCurriculumScore
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	if req.ClassCurriculumId == 0 {
		c.String(400, "课程Id不能为空")
		return
	}
	if req.ExaminationClassId == 0 {
		c.String(400, "考试Id不能为空")
		return
	}

	as, err := p.AchievementService.ListCurriculumScore(&req)
	if ut.IsValidateError(err) {
		c.String(400, err.Error())
		return
	}
	utee.Chk(err)

	c.JSON(200, utee.J{
		"data": as,
	})
}

func (p *Achievement) ListClassScore(c *gin.Context) {
	var req form.ListAchievementClassScore
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	if req.ExaminationClassId == 0 {
		c.String(400, "考试Id不能为空")
		return
	}

	as, err := p.AchievementService.ListClassScore(&req)
	if ut.IsValidateError(err) {
		c.String(400, err.Error())
		return
	}
	utee.Chk(err)

	c.JSON(200, utee.J{
		"data": as,
	})
}

func (p *Achievement) ListSearchScore(c *gin.Context) {
	var req form.ListAchievementSearchSource
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	if (req.StartTime != "" && req.EndTime == "") || (req.StartTime == "" && req.EndTime != "") {
		c.String(400, "开始和结束时间必须同时有值")
		return
	}
	as, total, err := p.AchievementService.ListSearchScore(&req)
	if ut.IsValidateError(err) {
		c.String(400, err.Error())
		return
	}
	utee.Chk(err)

	c.JSON(200, utee.J{
		"page":  req.Page,
		"limit": req.Limit,
		"total": total,
		"data":  as,
	})
}

func (p *Achievement) SaveScore(c *gin.Context) {
	var req []*form.SaveAchievementScore
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	val := func(req *form.SaveAchievementScore) (bool, string) {
		if req.ExaminationClassId == 0 {
			return false, "考试Id不能为空"
		}
		if req.ClassCurriculumId == 0 {
			return false, "班级课程Id不能为空"
		}
		if req.StudentId == 0 {
			return false, "学生Id不能为空"
		}
		return true, ""
	}
	for _, v := range req {
		ok, str := val(v)
		if !ok {
			c.String(400, str)
			return
		}
	}

	utee.Chk(p.AchievementService.SaveScore(req))

	c.String(200, "OK")
}

func (p *Achievement) Archived(c *gin.Context) {
	var req struct {
		ExaminationClassId int `json:"examination_class_id"`
	}
	if err := c.Bind(&req); err != nil {
		c.String(400, "参数错误")
		return
	}
	err := p.AchievementService.Archived(req.ExaminationClassId)
	if ut.IsValidateError(err) {
		c.String(400, err.Error())
		return
	}
	utee.Chk(err)
	c.String(200, "OK")
}
