package goc_mini_server

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xunull/goc/goc_mini_server/base_api"
	"net/http"
)

func (s *MiniServer) initMonitorRouter() {
	s.engine.GET("/metrics", transferHandler(promhttp.Handler()))
	s.engine.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
}

func (s *MiniServer) initSysRouter() {
	group := s.engine.Group("sys")
	{
		group.GET("routes", base_api.GetRoutes)
	}
}

func (s *MiniServer) initPublicGroup() {
	s.pubGroup = s.engine.Group("")
}

func (s *MiniServer) initStatic(path string) {
	if path != "" {
		s.engine.LoadHTMLGlob(path)
	}
}

func (s *MiniServer) initPrivateGroup() {
	s.priGroup = s.engine.Group("")
	{
		if s.Option.GinAccounts != nil {
			s.priGroup.Use(gin.BasicAuth(s.Option.GinAccounts))
		}
	}
}

func (s *MiniServer) init404() {
	s.engine.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
}

func transferHandler(handler http.Handler) gin.HandlerFunc {
	return func(context *gin.Context) {
		handler.ServeHTTP(context.Writer, context.Request)
	}
}
