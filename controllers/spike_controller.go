package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	R "shopping/response"
	"shopping/services"
	"shopping/utils"
)

type SpikeController struct {
	spikeService services.SpikeServiceImp
}

func (c *SpikeController) Shopping(ctx *gin.Context) {
	var spikeServiceUri services.SpikeServiceUri
	if err := ctx.ShouldBindUri(&spikeServiceUri); err == nil {
		userInfo, ok := ctx.Get("jwtUserInfo")
		if ok {
			R.Error(ctx, "系统错误", nil)
			return
		}
		info := userInfo.(utils.JwtUserInfo)
		if err := c.spikeService.Shopping(&info, spikeServiceUri.Id, ctx.GetHeader("Authorization")); err == nil {
			R.Ok(ctx, "抢购成功！", nil)
			return
		} else {
			R.Response(ctx, http.StatusNoContent, err.Error(), nil, http.StatusNoContent)
			return
		}
	} else {
		R.Response(ctx, http.StatusUnprocessableEntity, "参数错误", err.Error(), http.StatusUnprocessableEntity)
		return
	}
}
