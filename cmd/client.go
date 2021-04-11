package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/coder2z/g-saber/xcfg"
	"os"
	"shopping/models"
	"shopping/repositories"
	"shopping/services"
	"shopping/utils"
)

//消息队列客户端用于写入数据库防止暴库

var clientcfg string

func main() {

	flag.StringVar(&clientcfg, "c", "config/config.toml", "-c 	your config path")

	flag.Parse()

	utils.InitLog()

	file, err := os.Open(clientcfg)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = xcfg.LoadFromReader(file, toml.Unmarshal)

	if err != nil {
		panic(err)
	}

	models.Init()

	commodityRepository := repositories.CommodityRepository{Db: models.MysqlHandler}
	commodityService := services.CommodityService{CommodityRepository: &commodityRepository}

	orderRepository := repositories.OrderRepository{Db: models.MysqlHandler}
	orderService := services.OrderService{OrderRepository: &orderRepository}

	services.UpdateMQUrl(xcfg.GetString("mq.url"))

	simple := services.NewRabbitMQSimple("myxy99Shopping")
	simple.ConsumeSimple(&orderService, &commodityService)
}
