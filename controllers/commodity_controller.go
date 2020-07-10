package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	R "shopping/response"
	"shopping/services"
)

type CommodityControllerImp interface {
	GetCommodityById(ctx *gin.Context)
	GetCommodity(ctx *gin.Context)
	UpdateCommodity(ctx *gin.Context)
	AddCommodity(ctx *gin.Context)
	DelCommodity(ctx *gin.Context)
}

type CommodityController struct {
	CommodityService services.CommodityServiceImp `inject:""`
}

func (c *CommodityController) GetCommodityById(ctx *gin.Context) {
	var getCommodityIdService services.GetCommodityIdService
	if err := ctx.ShouldBindUri(&getCommodityIdService); err == nil {
		if data, err := c.CommodityService.GetCommodityById(&getCommodityIdService); err == nil {
			R.Ok(ctx, "成功", data)
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusUnprocessableEntity, "参数错误", err.Error(), http.StatusUnprocessableEntity)
	}
	return
}

func (c *CommodityController) GetCommodity(ctx *gin.Context) {
	var getCommodityPageService services.GetCommodityPageService
	if err := ctx.ShouldBind(&getCommodityPageService); err == nil {
		if data, err := c.CommodityService.GetCommodityPage(&getCommodityPageService); err == nil {
			R.Ok(ctx, "成功", data)
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusUnprocessableEntity, "参数错误", err.Error(), http.StatusUnprocessableEntity)
	}
	return
}

func (c *CommodityController) UpdateCommodity(ctx *gin.Context) {
	var (
		err                   error
		commodityFormService  services.CommodityFormService
		getCommodityIdService services.GetCommodityIdService
	)

	err = ctx.ShouldBind(&commodityFormService)

	err = ctx.ShouldBindUri(&getCommodityIdService)
	if err == nil {
		if err := c.CommodityService.UpdateCommodity(&commodityFormService, &getCommodityIdService); err == nil {
			R.Ok(ctx, "成功", nil)
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusUnprocessableEntity, "参数错误", err.Error(), http.StatusUnprocessableEntity)
	}
	return
}

func (c *CommodityController) AddCommodity(ctx *gin.Context) {
	var commodityFormService services.CommodityFormService
	if err := ctx.ShouldBind(&commodityFormService); err == nil {
		if err := c.CommodityService.AddCommodity(&commodityFormService); err == nil {
			R.Ok(ctx, "成功", nil)
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusUnprocessableEntity, "参数错误", err.Error(), http.StatusUnprocessableEntity)
	}
	return
}

func (c *CommodityController) DelCommodity(ctx *gin.Context) {
	var getCommodityIdService services.GetCommodityIdService
	if err := ctx.ShouldBindUri(&getCommodityIdService); err == nil {
		if err := c.CommodityService.DelCommodity(&getCommodityIdService); err == nil {
			R.Ok(ctx, "成功", nil)
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusUnprocessableEntity, "参数错误", err.Error(), http.StatusUnprocessableEntity)
	}
	return
}
