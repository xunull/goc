package extensions

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func InitPprofServer(addr string) error {
	router := gin.Default()
	pprof.Register(router)
	group := router.Group("/admin")
	{
		pprof.RouteRegister(group, "pprof")
	}
	err := router.Run(addr)
	return err
}
