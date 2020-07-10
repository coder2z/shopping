package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	R "shopping/response"
	"shopping/utils"
)

type Handler struct {
	Consistent utils.ConsistentHashImp
	LocalHost  string
	HostList   []string
	Port       string
}

func NewHandler(consistent utils.ConsistentHashImp, hostList []string, port string) *Handler {
	return &Handler{
		Consistent: consistent,
		HostList:   hostList,
		LocalHost:  "",
		Port:       port,
	}
}

func (h *Handler) Shopping(ctx *gin.Context) {
	//ctx.Redirect(http.StatusFound, "http://127.0.0.2:8081")
	userInfo, ok := ctx.Get("jwtUserInfo")
	if ok {
		R.Error(ctx, "系统错误", nil)
		return
	}
	info := userInfo.(utils.JwtUserInfo)
	fmt.Println(info)

	//id := strconv.Itoa(int(info.Id))
	//ip, err := h.Consistent.Get(id)
	//if err != nil {
	//	R.Error(ctx, err.Error(), nil)
	//	return
	//}
	//if ip == h.LocalHost {
	//	//本地处理
	//} else {
	//	//代理处理
	//}
}
