package goc_mini_server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/xunull/goc/commonx"
	"github.com/xunull/goc/gin_app"
	"github.com/xunull/goc/gin_extend"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
)

type MiniServer struct {
	Name     string
	engine   *gin.Engine
	server   *http.Server
	pubGroup *gin.RouterGroup
	priGroup *gin.RouterGroup
	Option   serverOption
}

func (s *MiniServer) init() {
	s.initBasicMiddleware()
	s.initBasicRouter()
}

func (s *MiniServer) initBasicMiddleware() {

	s.initLogger()
	s.initCors()
	s.initNoCache()
	//s.initBasicAuth()
}

func (s *MiniServer) initBasicRouter() {
	s.initMonitorRouter()
	s.initPublicGroup()
	s.initPrivateGroup()
	s.init404()
}

func (s *MiniServer) registerApp(app gin_app.App) {

	app.RegisterEngine(s.engine)

	pub := s.pubGroup.Group(strings.ToLower(s.Name))
	app.RegisterPubGroup(pub)

	pri := s.priGroup.Group(strings.ToLower(s.Name))
	app.RegisterPriGroup(pri)

	s.initStatic(app.GetHtmlRoot())
}

func (s *MiniServer) formatAddr() string {
	return fmt.Sprintf("%s:%d", s.Option.Host, s.Option.Port)
}

func (s *MiniServer) GetServerStatus() (*ServerStatus, error) {
	status, err := commonx.GetServerStatus()
	if err != nil {
		return nil, err
	}
	info := &ServerStatus{
		SystemStatus: status,
		Host:         s.Option.Host,
		Port:         s.Option.Port,
		RouteCount:   len(s.engine.Routes()),
		RouteTable:   s.GetRouteTable(),
	}
	return info, nil
}

func (s *MiniServer) GetRouteTable() string {
	return gin_extend.GetRouteTable(s.engine)
}

// ---------------------------------------------------------------------------------------------------------------------

func (s *MiniServer) Run() {
	if err := s.server.ListenAndServe(); err != nil {
		debug.PrintStack()
		log.Error().Err(err).Msg("mini server listen failed")
	}
}

func (s *MiniServer) Shutdown() {
	if err := s.server.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Msg("mini server shutdown failed")
	}
}

func (s *MiniServer) QuitWatch() {
	c := make(chan os.Signal)
	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-c
}
