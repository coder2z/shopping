package test

import (
	"fmt"
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"shopping/controllers"
	"shopping/middleware"
	"shopping/models"
	"shopping/repositories"
	"shopping/services"
	"testing"
)

type DBEngine struct {
	Name string
}

type UserDB struct {
	Db *DBEngine `inject:""`
}

type UserService struct {
	Db *UserDB `inject:""`
}

type App struct {
	Name string
	User *UserService `inject:""`
}

func (a *App) Create() string {
	return "create app, in db name:" + a.User.Db.Db.Name + " app name :" + a.Name
}

type Object struct {
	App *App
}

func Init() *Object {
	var g inject.Graph

	// 不适用依赖注入
	//a := DBEngine{Name: "db1"}
	//b := UserDB{&a}
	//c := UserService{&b}
	//app := App{Name: "go-app", User: &c}

	app := App{Name: "go-app"}

	_ = g.Provide(
		&inject.Object{Value: &DBEngine{Name: "db1"}},
		&inject.Object{Value: &app},
	)
	_ = g.Populate()

	return &Object{
		App: &app,
	}

}
func TestMains(t *testing.T) {
	obj := Init()
	fmt.Println(obj.App.Create())
}

func TestInject(t *testing.T) {
	models.Init()
	models.MysqlHandler.AutoMigrate(models.User{})
	//使用 Inject
	var userController controllers.UserController
	var injector inject.Graph
	_ = injector.Provide(
		&inject.Object{Value: &repositories.UserManagerRepository{Db: models.MysqlHandler}},
		&inject.Object{Value: &services.UserService{}},
		&inject.Object{Value: &userController},
	)
	_ = injector.Populate()

	app := gin.Default()
	api := app.Group("/api")
	{
		api.POST("/login", userController.Login)
		api.POST("/register", userController.Register)
		api.GET("/me", middleware.Auth(), userController.Info)
	}

	_ = app.Run(":8080")

}

func TestNoInject(t *testing.T) {

	//models.Init()
	//models.MysqlHandler.AutoMigrate(models.User{})
	////不使用 Inject
	//repository := repositories.NewUserRepository()
	//userServices := services.NewUserServices(repository)
	//controller := controllers.NewUserController(userServices)
	//
	//app := gin.Default()
	//api := app.Group("/api")
	//{
	//	api.POST("/login", controller.Login)
	//	api.POST("/register", controller.Register)
	//	api.GET("/me", middleware.Auth(), controller.Info)
	//}
	//
	//_ = app.Run(":8080")

}

func TestExample(t *testing.T) {

}
