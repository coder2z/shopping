package main

import (
	"shopping/models"
	"shopping/repositories"
	"shopping/services"
)

//消息队列客户端用于写入数据库防止暴库

func main() {
	models.Init()

	commodityRepository := repositories.CommodityRepository{Db: models.MysqlHandler}
	commodityService := services.CommodityService{CommodityRepository: &commodityRepository}

	orderRepository := repositories.OrderRepository{Db: models.MysqlHandler}
	orderService := services.OrderService{OrderRepository: &orderRepository}

	simple := services.NewRabbitMQSimple("myxy99Shopping")
	simple.ConsumeSimple(&orderService, &commodityService)
}
