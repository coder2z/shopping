package main

import (
	"shopping/models"
	"shopping/router"
	"shopping/utils"
)

func main() {
	utils.InitLog()

	models.Init()
	models.MysqlHandler.AutoMigrate(models.User{})
	models.MysqlHandler.AutoMigrate(models.Commodity{})

	app := router.InitRouter()

	_ = app.Run(":8080")
}
