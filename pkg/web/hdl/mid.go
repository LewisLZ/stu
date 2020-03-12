package hdl

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/quexer/utee"

	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/model"
)

type Mid struct {
	Ds *datasource.Ds
}

func (p *Mid) AuthRequiredMiddleware(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get(SessionUserUid)
	if v == nil {
		c.AbortWithStatus(401)
		return
	}
	id, ok := v.(int)
	if !ok {
		c.AbortWithStatus(401)
		return
	}

	var u model.User
	err := p.Ds.Db.First(&u, id).Error
	if datasource.NotFound(err) {
		c.AbortWithStatus(401)
		return
	}
	utee.Chk(err)

	passwd := session.Get(SessionUserPwt)
	if passwd == nil {
		c.AbortWithStatus(401)
		return
	}
	if passwd != u.Passwd {
		c.AbortWithStatus(401)
		return
	}

	// 每次认证通过都回写cookie, 保证操作期间cookie 不过期
	// 稍耗一点流量， 但不是问题
	err = SetWebSession(c, &u)
	utee.Chk(err)
	c.Set(CTX_USER_SESSION, &u)
}
