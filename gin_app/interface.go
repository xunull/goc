package gin_app

import "github.com/gin-gonic/gin"

type App interface {
	RegisterPriGroup(*gin.RouterGroup)
	RegisterPubGroup(*gin.RouterGroup)
	RegisterEngine(engine *gin.Engine)
	GetPriAuth()
	GetHtmlRoot() string
}
