package router

import (
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"log"
	"shopping/controllers"
	"shopping/middleware"
	"shopping/models"
	"shopping/repositories"
	"shopping/services"
)

func InitRouter() *gin.Engine {
	var userController controllers.UserController
	var commodityController controllers.CommodityController
	var orderController controllers.OrderController

	//依赖注入
	var injector inject.Graph
	err := injector.Provide(
		&inject.Object{Value: &repositories.UserManagerRepository{Db: models.MysqlHandler}},
		&inject.Object{Value: &services.UserService{}},
		&inject.Object{Value: &userController},

		&inject.Object{Value: &repositories.CommodityRepository{Db: models.MysqlHandler}},
		&inject.Object{Value: &services.CommodityService{}},
		&inject.Object{Value: &commodityController},

		&inject.Object{Value: &repositories.OrderRepository{Db: models.MysqlHandler}},
		&inject.Object{Value: &services.OrderService{}},
		&inject.Object{Value: &orderController},
	)
	if err != nil {
		log.Fatal("inject fatal: ", err)
	}
	if err := injector.Populate(); err != nil {
		log.Fatal("inject fatal: ", err)
	}

	//gin
	app := gin.Default()
	api := app.Group("/api")
	{
		api.POST("/login", userController.Login)
		api.POST("/register", userController.Register)
		api.GET("/me", middleware.Auth(), userController.Info)

		api.POST("/commodity", commodityController.AddCommodity)
		api.DELETE("/commodity/:id", commodityController.DelCommodity)
		api.GET("/commodity/:id", commodityController.GetCommodityById)
		api.GET("/commodity", commodityController.GetCommodity)
		api.PUT("/commodity/:id", commodityController.UpdateCommodity)

		api.GET("/order", orderController.Get)
	}
	return app
}
