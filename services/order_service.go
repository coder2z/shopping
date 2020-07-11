package services

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"shopping/models"
	"shopping/repositories"
	"shopping/utils"
	"strconv"
)

type GetOrderPageService struct {
	PageSize int `json:"page_size" form:"pageSize" binding:"required,numeric"`
	Page     int `json:"page" form:"page" binding:"required,numeric"`
}

type MessageService struct {
	models.Message
}

type OrderInfo struct {
	UserName       string `json:"user_name"`
	OrderId        string `json:"order_id"`
	UserEmail      string `json:"user_email"`
	Tel            string `json:"tel"`
	CommodityName  string `json:"commodity_name"`
	CommodityLink  string `json:"commodity_link"`
	CommodityPrice string `json:"commodity_price"`
}

type OrderServiceImp interface {
	GetOrder(*GetOrderPageService) (*utils.Page, error)
	Add(*MessageService) error
}

type OrderService struct {
	OrderRepository repositories.OrderRepositoryImp `inject:""`
}

func (s *OrderService) Add(m *MessageService) (err error) {
	node, err := utils.NewWorker(int64(m.UserID))
	if err != nil {
		utils.Log.WithFields(log.Fields{"errMsg": err.Error()}).Warningln("创建订单失败")
		return
	}
	orderId := strconv.FormatInt(node.GetId(), 10)
	orderInfo := models.Order{
		UserId:      m.UserID,
		CommodityId: m.CommodityId,
		OrderId:     orderId,
	}
	err = s.OrderRepository.Add(&orderInfo)
	if err != nil {
		utils.Log.WithFields(log.Fields{"errMsg": err.Error()}).Warningln("创建订单失败")
		return
	}
	return
}

func (s *OrderService) GetOrder(pageInfo *GetOrderPageService) (p *utils.Page, err error) {
	o, total, err := s.OrderRepository.GetSize(pageInfo.PageSize*(pageInfo.Page-1), pageInfo.PageSize)
	if err != nil {
		utils.Log.WithFields(log.Fields{"errMsg": err.Error()}).Warningln("获取数据失败")
		return p, errors.New("获取数据失败！")
	}
	var list []OrderInfo
	var tmpV OrderInfo
	for _, V := range *o {
		tmpV = OrderInfo{
			UserName:       V.User.UserName,
			OrderId:        V.OrderId,
			UserEmail:      V.User.Email,
			Tel:            V.User.Tel,
			CommodityName:  V.Commodity.Name,
			CommodityLink:  V.Commodity.Link,
			CommodityPrice: V.Commodity.Price,
		}
		list = append(list, tmpV)
	}
	pages := utils.PageUtil(total, pageInfo.Page, pageInfo.PageSize, list)
	return &pages, err
}
