package services

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"shopping/models"
	"shopping/repositories"
	"shopping/utils"
)

type CommodityFormService struct {
	Name      string `json:"name" form:"name" binding:"required"`
	Link      string `json:"link" form:"link" binding:"required"`
	Price     string `json:"price" form:"price" binding:"required"`
	Stock     int    `json:"stock" form:"stock" binding:"required,numeric"`
	StartTime int64 `json:"startTime" form:"start_time" binding:"required"`
}

type GetCommodityPageService struct {
	PageSize int `json:"page_size" form:"pageSize" binding:"required,numeric"`
	Page     int `json:"page" form:"page" binding:"required,numeric"`
}

type GetCommodityIdService struct {
	Id int `json:"id" uri:"id" binding:"required,numeric"`
}

type CommodityServiceImp interface {
	//获取商品信息
	GetCommodityById(*GetCommodityIdService) (*models.Commodity, error)
	//获取全部商品
	GetCommodityAll() (*[]models.Commodity, error)
	//获取分页商品
	GetCommodityPage(*GetCommodityPageService) (*utils.Page, error)
	//更新商品
	UpdateCommodity(*CommodityFormService, *GetCommodityIdService) error
	//添加商品
	AddCommodity(*CommodityFormService) error
	//删除商品
	DelCommodity(*GetCommodityIdService) error
	//商品库存减去一个
	SubNumberOne(int) error
}

type CommodityService struct {
	CommodityRepository repositories.CommodityRepositoryImp `inject:""`
}

func (s *CommodityService) GetCommodityById(idForm *GetCommodityIdService) (c *models.Commodity, err error) {
	c, err = s.CommodityRepository.GetById(idForm.Id)
	if err != nil {
		utils.Log.WithFields(log.Fields{"errMsg": err.Error()}).Warningln("获取数据失败")
		return c, errors.New("获取数据失败！")
	}
	return
}

func (s *CommodityService) GetCommodityAll() (c *[]models.Commodity, err error) {
	c, err = s.CommodityRepository.GetAll()
	if err != nil {
		utils.Log.WithFields(log.Fields{"errMsg": err.Error()}).Warningln("获取数据失败")
		return c, errors.New("获取数据失败！")
	}
	return
}

func (s *CommodityService) GetCommodityPage(pageInfo *GetCommodityPageService) (p *utils.Page, err error) {
	c, total, err := s.CommodityRepository.GetSize((pageInfo.Page-1)*pageInfo.PageSize, pageInfo.PageSize)
	if err != nil {
		utils.Log.WithFields(log.Fields{"errMsg": err.Error()}).Warningln("获取数据失败")
		return p, errors.New("获取数据失败！")
	}
	pages := utils.PageUtil(total, pageInfo.Page, pageInfo.PageSize, c)
	return &pages, err
}

func (s *CommodityService) UpdateCommodity(updateInfo *CommodityFormService, idForm *GetCommodityIdService) (err error) {
	info := &models.Commodity{
		Name:      updateInfo.Name,
		Link:      updateInfo.Link,
		Price:     updateInfo.Price,
		Stock:     updateInfo.Stock,
		StartTime: updateInfo.StartTime}
	err = s.CommodityRepository.Update(idForm.Id, info)
	if err != nil {
		utils.Log.WithFields(log.Fields{"errMsg": err.Error()}).Warningln("修改商品数据失败")
		return errors.New("修改商品数据失败！")
	}
	return
}

func (s *CommodityService) AddCommodity(addInfo *CommodityFormService) (err error) {
	info := &models.Commodity{
		Name:      addInfo.Name,
		Link:      addInfo.Link,
		Price:     addInfo.Price,
		Stock:     addInfo.Stock,
		StartTime: addInfo.StartTime}
	err = s.CommodityRepository.Add(info)
	if err != nil {
		utils.Log.WithFields(log.Fields{"errMsg": err.Error()}).Warningln("添加商品数据失败")
		return errors.New("添加商品数据失败！")
	}
	return
}

func (s *CommodityService) DelCommodity(idForm *GetCommodityIdService) (err error) {
	err = s.CommodityRepository.Del(idForm.Id)
	if err != nil {
		utils.Log.WithFields(log.Fields{"errMsg": err.Error()}).Warningln("删除商品数据失败")
		return errors.New("删除商品数据失败！")
	}
	return
}

func (s *CommodityService) SubNumberOne(commodityId int) (err error) {
	err = s.CommodityRepository.UpdateStockMinusOne(commodityId)
	if err != nil {
		utils.Log.WithFields(log.Fields{"errMsg": err.Error()}).Warningln("商品库存修改失败")
		return errors.New("商品库存修改失败")
	}
	return
}
