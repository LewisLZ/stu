package hdl

import (
	"github.com/gin-gonic/gin"

	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/ut"
)

const (
	SessionUserUid   = "user_uid"
	SessionUserPwt   = "user_pwt"
	CTX_USER_SESSION = "ctx_user_session"
)

func MustGetUser(c *gin.Context) *model.User {
	return c.MustGet(CTX_USER_SESSION).(*model.User)
}

func SetWebSession(c *gin.Context, user *model.User) error {
	sessionKeys := []ut.SessionKey{
		{
			Key:   SessionUserUid,
			Value: user.Id,
		}, {
			Key:   SessionUserPwt,
			Value: user.Passwd,
		},
	}

	return ut.SetWebSession(c, ut.SessionOpt{
		Keys: sessionKeys,
	})
}

func DelWebSession(c *gin.Context) error {
	sessionKeys := []ut.SessionKey{
		{
			Key: SessionUserUid,
		},
	}
	return ut.DelWebSession(c, ut.SessionOpt{
		Keys: sessionKeys,
	})
}
