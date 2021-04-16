package goc_mini_server

import (
	"github.com/gin-gonic/gin"
	"github.com/xunull/goc/gin_app"
	"github.com/xunull/goc/goc_mini_server/base_api"
	"net/http"
	"time"
)

type serverOption struct {
	Host        string
	Port        int
	GinAccounts gin.Accounts
	LogLevel    string
	LogDir      string
	ForceAuth   bool
}

type Option func(o *serverOption)

func WithForceAuth() Option {
	return func(o *serverOption) {
		o.ForceAuth = true
	}
}

func WithLogDir(path string) Option {
	return func(o *serverOption) {
		o.LogDir = path
	}
}

func WithLogLevel(level string) Option {
	return func(o *serverOption) {
		o.LogLevel = level
	}
}

func WithPort(port int) Option {
	return func(o *serverOption) {
		o.Port = port
	}
}

func WithBasicAuth(userPass map[string]string) Option {
	return func(o *serverOption) {
		o.GinAccounts = userPass
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func RunApp(name string, app gin_app.App, opts ...Option) *MiniServer {
	d := &serverOption{
		Host:      "127.0.0.1",
		Port:      0,
		LogLevel:  "debug",
		ForceAuth: false,
	}
	for _, o := range opts {
		o(d)
	}

	ms := MiniServer{Option: *d,
		Name: name}

	engine := gin.Default()
	ms.engine = engine

	base_api.Engine = engine

	server := &http.Server{
		Addr:           ms.formatAddr(),
		Handler:        ms.engine,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	ms.server = server
	ms.init()

	ms.registerApp(app)

	go ms.Run()

	return &ms

}
