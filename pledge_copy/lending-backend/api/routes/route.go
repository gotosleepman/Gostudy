package routes

import (
	"lending-copy/api/controllers"
	"lending-copy/config"

	"github.com/gin-gonic/gin"
)

func InitRoute(e *gin.Engine) *gin.Engine {
	v1 := e.Group("/api/v" + config.Config.Env.Version)
	poolController := controllers.PoolController{}
	v1.GET("/poolBaseInfo", poolController.PoolBaseInfo)
	v1.GET("/poolDataInfo", poolController.PoolDataInfo)
	v1.GET("/token", poolController.TokenList)
	v1.POST("/pool/search", poolController.Search)
	return e
}
