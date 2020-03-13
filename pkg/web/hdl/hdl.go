package hdl

import (
	"github.com/gin-gonic/gin"

	"liuyu/stu/pkg/datasource"
)

type Hdl struct {
	Ds         *datasource.Ds
	Student    *Student
	Teacher    *Teacher
	Pub        *Pub
	User       *User
	Mid        *Mid
	Curriculum *Curriculum
	Class      *Class
}

func (p *Hdl) Mount(rg *gin.RouterGroup) {
	p.Pub.Mount(rg.Group("/pub"))

	sys := rg.Group("/sys", p.Mid.AuthRequiredMiddleware)
	{
		p.User.Mount(sys.Group("/user"))
		p.Teacher.Mount(sys.Group("/teacher"))
		p.Student.Mount(sys.Group("/student"))
		p.Curriculum.Mount(sys.Group("/curriculum"))
		p.Class.Mount(sys.Group("/class"))
	}
}
