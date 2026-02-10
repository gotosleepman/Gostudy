package main

import (
	"myblog-backend/config"
	"myblog-backend/controllers"
	"myblog-backend/database"
	"myblog-backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化日志
	middleware.InitLogger()

	// 连接数据库
	database.ConnectDatabase()

	// 初始化Gin
	if cfg.LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 使用中间件
	r.Use(middleware.LoggerMiddleware())
	r.Use(gin.Recovery()) // 恢复panic

	// 初始化控制器
	authController := &controllers.AuthController{}
	postController := &controllers.PostController{}
	commentController := &controllers.CommentController{}

	// 健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "myblog-backend",
		})
	})

	// API路由组
	api := r.Group("/api/v1")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.GET("/profile", middleware.AuthMiddleware(), authController.GetProfile)
		}

		// 文章路由
		posts := api.Group("/posts")
		{
			posts.GET("", middleware.OptionalAuthMiddleware(), postController.GetPosts)
			posts.GET("/:id", middleware.OptionalAuthMiddleware(), postController.GetPost)
			posts.POST("", middleware.AuthMiddleware(), postController.CreatePost)
			posts.PUT("/:id", middleware.AuthMiddleware(), postController.UpdatePost)
			posts.DELETE("/:id", middleware.AuthMiddleware(), postController.DeletePost)
		}

		// 评论路由
		comments := api.Group("/comments")
		{
			comments.GET("/post/:postId", commentController.GetComments)
			comments.POST("", middleware.AuthMiddleware(), commentController.CreateComment)
			comments.DELETE("/:id", middleware.AuthMiddleware(), commentController.DeleteComment)
		}
	}

	// 启动服务器
	r.Run(":" + cfg.ServerPort)
}
