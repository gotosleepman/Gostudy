package middleware

import (
	"fmt"
	"myblog-backend/config"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var logger *zap.Logger

// InitLogger 初始化日志记录器
func InitLogger() {
	cfg := config.LoadConfig()

	var err error
	if cfg.LogLevel == "debug" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	defer logger.Sync()
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求
		c.Next()

		// 记录日志
		end := time.Now()
		latency := end.Sub(start)

		logger.Info("Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}

// ErrorLogger 错误日志记录
func ErrorLogger(err error, context string) {
	if logger != nil {
		logger.Error(context,
			zap.Error(err),
			zap.String("time", time.Now().Format(time.RFC3339)),
		)
	}
}
