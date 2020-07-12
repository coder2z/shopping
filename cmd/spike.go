package main

import (
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"shopping/controllers"
	"shopping/middleware"
	"shopping/models"
	"shopping/repositories"
	R "shopping/response"
	"shopping/services"
	"shopping/utils"
	"strconv"
	"sync"
)

//基于hash环的分布式秒杀
var (
	//分布式集群地址
	hostList = []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}
	//端口
	port = "8081"
	//记录现在的秒杀商品的数量
	commodityCache map[int]models.Commodity
	//锁
	mutex sync.Mutex

	//hash环
	consistent utils.ConsistentHashImp
)

func main() {
	consistent = utils.NewConsistent(20)
	for _, v := range hostList {
		consistent.Add(v)
	}
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
	for _, value := range *commodityList {
		commodityCache[int(value.ID)] = value
	}

	app := gin.Default()
	ip, err := utils.GetIp()
	if err != nil {
		utils.Log.WithFields(log.Fields{"errMsg": err.Error()}).Panic("ip获取失败")
		os.Exit(1)
		return
	}
	//ip = "127.0.0.3"
	simple := services.NewRabbitMQSimple("myxy99Shopping")
	spikeService := &services.SpikeService{
		Consistent:       consistent,
		LocalHost:        ip,
		HostList:         hostList,
		Port:             port,
		CommodityCache:   &commodityCache,
		RabbitMqValidate: simple,
	}

	spikeController := &controllers.SpikeController{SpikeService: spikeService} //, middleware.Auth()

	limiter := tollbooth.NewLimiter(1, nil)
	app.GET("/:uid/spike/:id", tollbooth_gin.LimitHandler(limiter), Ip(consistent, ip), middleware.Auth(), spikeController.Shopping)

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
				//代理处理
				res, _, _ := utils.GetCurl(fmt.Sprintf("http://%v:%v/spike/%v", ip, port, c.Param("id")), c.GetHeader("Authorization"))
				if res.StatusCode == 200 {
					mutex.Lock()
					defer mutex.Unlock()
					commodityCache[spikeServiceUri.Id].Stock--
					R.Response(c, http.StatusCreated, "成功抢到", nil, http.StatusCreated)
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
