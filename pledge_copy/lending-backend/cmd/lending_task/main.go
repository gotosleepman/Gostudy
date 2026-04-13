package main

import (
	"lending-copy/db"
	"lending-copy/schedule/models"
	"lending-copy/schedule/tasks"
)

func main() {
	db.InitMysql()
	db.InitRedis()
	models.InitTable()
	tasks.Task()
}
