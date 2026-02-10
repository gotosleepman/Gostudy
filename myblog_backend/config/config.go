package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort   string
	DatabaseURL  string
	JWTSecret    string
	JWTExpire    time.Duration
	LogLevel     string
	DatabaseType string // "mysql" or "sqlite"
}

func LoadConfig() *Config {
	// 从环境变量加载配置，如果没有设置则使用默认值
	serverPort := getEnv("SERVER_PORT", "8080")
	databaseURL := getEnv("DATABASE_URL", "blog.db")
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	jwtExpireStr := getEnv("JWT_EXPIRE_HOURS", "24")
	databaseType := getEnv("DATABASE_TYPE", "sqlite")
	logLevel := getEnv("LOG_LEVEL", "debug")

	jwtExpireHours, err := strconv.Atoi(jwtExpireStr)
	if err != nil {
		log.Fatal("Invalid JWT_EXPIRE_HOURS value")
	}

	// MySQL 连接示例：user:password@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local
	// SQLite 连接示例：blog.db

	return &Config{
		ServerPort:   serverPort,
		DatabaseURL:  databaseURL,
		JWTSecret:    jwtSecret,
		JWTExpire:    time.Duration(jwtExpireHours) * time.Hour,
		LogLevel:     logLevel,
		DatabaseType: databaseType,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
