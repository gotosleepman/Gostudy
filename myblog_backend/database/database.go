package database

import (
	"fmt"
	"log"
	"myblog-backend/config"
	"myblog-backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDatabase 连接数据库
func ConnectDatabase() {
	cfg := config.LoadConfig()
	var err error
	var dialect gorm.Dialector

	switch cfg.DatabaseType {
	case "mysql":
		dialect = mysql.Open(cfg.DatabaseURL)
	case "sqlite":
		dialect = sqlite.Open(cfg.DatabaseURL)
	default:
		log.Fatal("Unsupported database type")
	}

	DB, err = gorm.Open(dialect, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Database connection established")

	// 自动迁移
	err = DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database migrated successfully")
}
