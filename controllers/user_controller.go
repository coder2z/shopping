package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	R "shopping/response"
	"shopping/services"
	"shopping/utils"
)

type UserController struct {
	UserServices services.UserServiceImp `inject:""`
}

type UserImp interface {
	Login(*gin.Context)
	Register(*gin.Context)
	Info(*gin.Context)
}

//func NewUserController(userServices services.UserServiceImp) UserImp {
//	return &UserController{UserServices: userServices}
//}

func (l *UserController) Login(c *gin.Context) {
	var loginService services.UserLoginService
	if err := c.ShouldBind(&loginService); err == nil {
		if token, err := l.UserServices.Login(&loginService); err == nil {
			R.Ok(c, "成功", gin.H{
				"token": token,
			})
		} else {
			R.Error(c, err.Error(), nil)
		}
	} else {
		R.Response(c, http.StatusUnprocessableEntity, "参数错误", err.Error(), http.StatusUnprocessableEntity)
	}
}

func (l *UserController) Register(c *gin.Context) {
	var registerService services.UserRegisterService
	if err := c.ShouldBind(&registerService); err == nil {
		if err := l.UserServices.Register(&registerService); err == nil {
			R.Ok(c, "成功", nil)
		} else {
			R.Error(c, err.Error(), nil)
		}
	} else {
		R.Response(c, http.StatusUnprocessableEntity, "参数错误", err.Error(), http.StatusUnprocessableEntity)
	}
	return
}

func (l *UserController) Info(c *gin.Context) {
	userInfo, exists := c.Get("jwtUserInfo")
	if exists {
		info := userInfo.(utils.JwtUserInfo)
		if userInfo, err := l.UserServices.Info(info.Email); err == nil {
			R.Ok(c, "成功", gin.H{
				"email":    userInfo.Email,
				"userName": userInfo.UserName,
				"tel":      userInfo.Tel,
				"id":       userInfo.ID,
			})
		} else {
			R.Error(c, err.Error(), nil)
		}
	} else {
		R.Error(c, "失败", nil)
	}
}
