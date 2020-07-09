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

	//依赖注入
	var injector inject.Graph
	err := injector.Provide(
		&inject.Object{Value: &repositories.UserManagerRepository{Db: models.MysqlHandler}},
		&inject.Object{Value: &services.UserService{}},
		&inject.Object{Value: &userController},
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
	}
	return app
}
