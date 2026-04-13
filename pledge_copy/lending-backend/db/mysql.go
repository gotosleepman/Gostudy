package db

import (
	"fmt"
	"time"

	"lending-copy/config"
	"lending-copy/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitMysql() {
	mysqlConf := config.Config.Mysql
	log.Logger.Info("Init Mysql")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConf.UserName,
		mysqlConf.Password,
		mysqlConf.Address,
		mysqlConf.Port,
		mysqlConf.DbName)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Logger.Panic(fmt.Sprintf("mysql connention error ==>  %+v", err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Logger.Error("db.DB() err:" + err.Error())
	}
	sqlDB.SetMaxIdleConns(mysqlConf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(mysqlConf.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(mysqlConf.MaxLifeTime) * time.Second)
	Mysql = db
}
