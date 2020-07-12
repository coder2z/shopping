package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	R "shopping/response"
	"shopping/services"
	"shopping/utils"
)

type SpikeController struct {
	SpikeService services.SpikeServiceImp
}

func (c *SpikeController) Shopping(ctx *gin.Context) {
	s, ok := ctx.Get("spikeServiceUri")
	spikeServiceUri := s.(services.SpikeServiceUri)
	if ok {
		userInfo, ok := ctx.Get("jwtUserInfo")
		if !ok {
			R.Error(ctx, "系统错误", nil)
			return
		}
		info := userInfo.(utils.JwtUserInfo)
		if err := c.SpikeService.Shopping(&info, spikeServiceUri.Id, ctx.GetHeader("Authorization")); err == nil {
			R.Ok(ctx, "抢购成功！", nil)
			return
		} else {
			R.Response(ctx, http.StatusCreated, err.Error(), nil, http.StatusCreated)
			return
		}
	} else {
		R.Response(ctx, http.StatusUnprocessableEntity, "参数错误", nil, http.StatusUnprocessableEntity)
		return
	}
}
