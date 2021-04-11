package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/coder2z/g-saber/xcast"
	"github.com/coder2z/g-saber/xcfg"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"net/http"
	"os"
	"shopping/cache"
	"shopping/constant"
	"shopping/controllers"
	"shopping/discovery"
	"shopping/middleware"
	"shopping/models"
	"shopping/repositories"
	R "shopping/response"
	"shopping/services"
	"shopping/utils"
	"strconv"
	"time"
)

//基于hash环的分布式秒杀
var (
	//分布式集群地址
	hostList = []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}
	//端口
	port = "8081"

	//hash环
	consistent utils.ConsistentHashImp

 	scfg string
)

func main() {
	flag.StringVar(&scfg, "c", "config/config.toml", "-c 	your config path")

	flag.Parse()

	utils.InitLog()

	file, err := os.Open(scfg)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = xcfg.LoadFromReader(file, toml.Unmarshal)

	if err != nil {
		panic(err)
	}

	port = xcfg.GetString("server.port")

	ch := make(chan []discovery.ServerInfo)

	discovery.New(clientv3.Config{
		Endpoints: xcfg.GetStringSlice("etcd.endpoints"),
	}, ch)

	select {
	case x := <-ch:
		UpdateAddress(x)
	case <-time.After(time.Minute):
		panic("not resolve success in one minute")
	}

	go func() {
		for i := range ch {
			UpdateAddress(i)
		}
	}()

	cache.RedisHandle()
	//缓存所有需要秒杀的商品的库存
	models.Init()
	models.MysqlHandler.AutoMigrate(models.Order{})
	repository := &repositories.CommodityRepository{Db: models.MysqlHandler}
	service := &services.CommodityService{CommodityRepository: repository}
	commodityList, err := service.GetCommodityAll()

	if err != nil {
		utils.Log.WithFields(log.Fields{"errMsg": err.Error()}).Panic("缓存所有需要秒杀的商品的库存，获取库存失败")
		os.Exit(1)
		return
	}

	if commodityList == nil {
		utils.Log.Panic("无秒杀商品")
	}

	for _, value := range *commodityList {
		err = utils.AddStock(context.Background(),
			constant.SpikeKey.Format(value.ID),
			xcast.ToInt64(value.Stock))
		if err != nil {
			utils.Log.Panic(err)
		}
	}

	app := gin.Default()
	ip, err := utils.GetIp()
	if err != nil {
		utils.Log.WithFields(log.Fields{"errMsg": err.Error()}).Panic("ip获取失败")
		os.Exit(1)
		return
	}

	services.UpdateMQUrl(xcfg.GetString("mq.url"))

	simple := services.NewRabbitMQSimple("myxy99Shopping")
	spikeService := &services.SpikeService{
		RabbitMqValidate: simple,
	}

	spikeController := &controllers.SpikeController{SpikeService: spikeService}

	limiter := tollbooth.NewLimiter(1, nil)
	app.GET("/:uid/spike/:id", tollbooth_gin.LimitHandler(limiter), middleware.Auth(), Ip(consistent, ip), spikeController.Shopping)

	app.GET("/local/:uid/spike/:id", spikeController.Shopping)

	app.GET("/", tollbooth_gin.LimitHandler(limiter), func(context *gin.Context) {
		context.JSON(200, gin.H{"data": 1})
	})

	_ = app.Run(fmt.Sprintf(":%v", port))
}
func Ip(Consistent utils.ConsistentHashImp, LocalHost string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var spikeServiceUri services.SpikeServiceUri
		if err := c.ShouldBindUri(&spikeServiceUri); err == nil {
			c.Set("spikeServiceUri", spikeServiceUri)
			id := strconv.Itoa(spikeServiceUri.UId)
			ip, err := Consistent.Get(id)
			if err != nil {
				utils.Log.WithFields(log.Fields{"errMsg": err.Error()}).Warningln("hash环获取数据错误")
				R.Response(c, http.StatusInternalServerError, "服务器错误", err.Error(), http.StatusInternalServerError)
				c.Abort()
				return
			}
			if ip == LocalHost {
				c.Next()
				return
			} else {
				res, _, _ := utils.GetCurl(fmt.Sprintf("http://%s:%s/local/%s/spike/%s", ip, port, c.Param("uid"), c.Param("id")), c.GetHeader("Authorization"))
				if res.StatusCode == 200 {
					R.Response(c, http.StatusOK, "成功抢到", nil, http.StatusOK)
					c.Abort()
					return
				} else {
					R.Response(c, http.StatusCreated, "未抢到", nil, http.StatusCreated)
					c.Abort()
					return
				}
			}
		} else {
			R.Response(c, http.StatusUnprocessableEntity, "参数错误", err.Error(), http.StatusUnprocessableEntity)
			c.Abort()
			return
		}
	}
}

func UpdateAddress(info []discovery.ServerInfo) {
	var list []string
	for _, serverInfo := range info {
		list = append(list, serverInfo.Ip)
	}
	hostList = list

	consistent = utils.NewConsistent(20)
	for _, v := range hostList {
		consistent.Add(v)
	}
}
