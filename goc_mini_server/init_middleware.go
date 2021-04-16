package goc_mini_server

import (
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/xunull/goc/commonx"
	"github.com/xunull/goc/logx"
	"net/http"
	"time"
)

func (s *MiniServer) initNoCache() {
	noCache := func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		c.Next()
	}
	s.engine.Use(noCache)
}

func (s *MiniServer) initBasicAuth() {
	if s.Option.ForceAuth {
		engine := s.engine
		if s.Option.GinAccounts != nil {
			engine.Use(gin.BasicAuth(s.Option.GinAccounts))
		}
	}
}

func (s *MiniServer) initCors() {

	cc := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	cc.AllowAllOrigins = true
	s.engine.Use(cors.New(cc))
}

func (s *MiniServer) initLogger() {
	name := "gin-" + s.Name
	logger, err := logx.InitZapFileLogger(s.Option.LogLevel, s.Option.LogDir, name)
	if err != nil {
		commonx.CheckErrOrFatal(err)
	}

	s.engine.Use(ginzap.Ginzap(logger, time.RFC3339, true))
}
