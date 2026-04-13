package main

import (
	"lending-copy/api/middlewares"
	"lending-copy/api/routes"
	"lending-copy/api/static"
	"lending-copy/api/validate"
	"lending-copy/config"
	"lending-copy/db"
	schedmodels "lending-copy/schedule/models"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitMysql()
	db.InitRedis()
	schedmodels.InitTable()

	validate.BindingValidator()

	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	staticPath := static.GetCurrentAbPathByCaller()
	app.Static("/storage/", staticPath)
	app.Use(middlewares.Cors())
	routes.InitRoute(app)
	_ = app.Run(":" + config.Config.Env.Port)
}
