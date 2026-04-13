package models

import "lending-copy/db"

func InitTable() {
	db.Mysql.AutoMigrate(&PoolBase{})
	db.Mysql.AutoMigrate(&PoolData{})
	db.Mysql.AutoMigrate(&TokenInfo{})
}
