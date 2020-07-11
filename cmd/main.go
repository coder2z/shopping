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
	models.MysqlHandler.AutoMigrate(models.Order{})

	models.MysqlHandler.Model(&models.Order{}).AddForeignKey("user_id", "user(id)", "RESTRICT", "RESTRICT")
	models.MysqlHandler.Model(&models.Order{}).AddForeignKey("commodity_id", "commodity(id)", "RESTRICT", "RESTRICT")

	app := router.InitRouter()

	_ = app.Run(":8080")
}
