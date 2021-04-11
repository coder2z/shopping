package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/coder2z/g-saber/xcfg"
	"os"
	"shopping/models"
	"shopping/router"
	"shopping/utils"
)

var maincfg string

func main() {
	flag.StringVar(&maincfg, "c", "config/config.toml", "-c 	your config path")

	flag.Parse()

	utils.InitLog()

	file, err := os.Open(maincfg)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = xcfg.LoadFromReader(file, toml.Unmarshal)

	if err != nil {
		panic(err)
	}

	models.Init()
	models.MysqlHandler.AutoMigrate(models.User{})
	models.MysqlHandler.AutoMigrate(models.Commodity{})
	models.MysqlHandler.AutoMigrate(models.Order{})

	models.MysqlHandler.Model(&models.Order{}).AddForeignKey("user_id", "user(id)", "RESTRICT", "RESTRICT")
	models.MysqlHandler.Model(&models.Order{}).AddForeignKey("commodity_id", "commodity(id)", "RESTRICT", "RESTRICT")

	app := router.InitRouter()

	_ = app.Run(fmt.Sprintf(":%d", xcfg.GetInt("server.port")))
}
