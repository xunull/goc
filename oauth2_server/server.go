package oauth2_server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xunull/goc/commonx"
	"github.com/xunull/goc/goc_mini_server"
	"net/http"
	"net/url"
)

type Oauth2Server struct {
	Server             *goc_mini_server.MiniServer
	CallbackValuesChan chan url.Values
}

func (s *Oauth2Server) RegisterEngine(engine *gin.Engine) {

}

func (s *Oauth2Server) RegisterPriGroup(g *gin.RouterGroup) {

}

func (s *Oauth2Server) GetPriAuth() {

}

func (s *Oauth2Server) GetHtmlRoot() string {
	return ""
}

// ---------------------------------------------------------------------------------------------------------------------

func (s *Oauth2Server) RegisterPubGroup(g *gin.RouterGroup) {
	g.GET("auth_callback", func(c *gin.Context) {
		fmt.Println(c.Request.URL)
		s.CallbackValuesChan <- c.Request.URL.Query()
		c.Status(http.StatusOK)
	})
}

func (s *Oauth2Server) GetAuthCallbackUrl() string {

	return fmt.Sprintf("http://%s:%d/%s/%s",
		s.Server.Option.Host,
		s.Server.Option.Port,
		s.Server.Name,
		"auth_callback")
}

func (s *Oauth2Server) StartServer(opts ...goc_mini_server.Option) {

	server := goc_mini_server.RunApp("helper", s, opts...)
	status, err := server.GetServerStatus()
	commonx.CheckErrOrFatal(err)
	status.Output()
	s.Server = server

	server.QuitWatch()
}
