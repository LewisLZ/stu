package web

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/static"

	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/ut"
	"liuyu/stu/pkg/web/hdl"
)

type Web struct {
	Hdl *hdl.Hdl
	Ds  *datasource.Ds
}

func MidHSTS() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("X-Forwarded-Proto") == "https" {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}
	}
}

const (
	SESSION_SALT = "!You.pIN*DP[IPnqX9d2"
	SESSION_NAME = "web-session"
)

func (p *Web) CreateWebServer(webOpt *ut.WebOpt) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(MidHSTS())

	// Middlewares
	store := cookie.NewStore([]byte(SESSION_SALT))
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	})
	router.Use(sessions.Sessions(SESSION_NAME, store))

	router.Use(func(c *gin.Context) {
		c.Set(ut.CtxDS, p.Ds)
		c.Next()
	})

	router.Use(gin.Recovery(), static.Serve("/", static.LocalFile(webOpt.War, true)))

	router.GET("/bv", func(c *gin.Context) { c.String(200, fmt.Sprintf("binary version: %v\n", Version)) })

	api := router.Group("/api")
	{
		p.Hdl.Mount(api)
	}

	http.Handle("/", router)
	fmt.Printf("web serve on %s", webOpt.Addr)
	panic(http.ListenAndServe(webOpt.Addr, nil))
}
