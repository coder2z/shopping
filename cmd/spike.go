package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"shopping/controllers"
	"shopping/utils"
)

//基于hash环的分布式权限控制

var (
	//分布式集群地址
	hostList = []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}
	//端口
	port = "8081"
	//
	commoditys = map[int]int{}
	//hash环
	consistent utils.ConsistentHashImp
)

func main() {
	consistent = utils.NewConsistent(20)
	for _, v := range hostList {
		consistent.Add(v)
	}
	//缓存所有需要秒杀的商品的库存

	//
	app := gin.Default()
	handler := controllers.NewHandler(consistent, hostList, port) //, middleware.Auth()
	app.GET("/spike/:commodityId", handler.Shopping)
	app.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{"data": 1})
	})

	_ = app.Run(fmt.Sprintf(":%v", port))
}
