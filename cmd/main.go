package main

import (
	"shopping/models"
	"shopping/router"
)

func main() {
	models.Init()
	models.MysqlHandler.AutoMigrate(models.User{})

	app := router.InitRouter()

	_ = app.Run(":8080")
}
