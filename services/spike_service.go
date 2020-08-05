package services

import (
	"encoding/json"
	"errors"
	"shopping/models"
	"shopping/utils"
	"time"
)

type SpikeServiceUri struct {
	Id  int `uri:"id" json:"id" binding:"required,numeric"`
	UId int `uri:"uid" json:"uid" binding:"required,numeric"`
}

type SpikeServiceImp interface {
	Shopping(*utils.JwtUserInfo, int) error
}

type SpikeService struct {
	CommodityCache   *map[int]models.Commodity
	RabbitMqValidate *RabbitMQ
}

func (s *SpikeService) Shopping(info *utils.JwtUserInfo, commodityId int) (err error) {
	a := *s.CommodityCache
	if a[commodityId].StartTime > time.Now().Unix() {
		return errors.New("商品未开卖！")
	}
	//操作
	message := MessageService{
		models.Message{
			CommodityId: uint(commodityId),
			UserID:      uint(info.Id),
		},
	}
	byteMessage, err := json.Marshal(message)
	if err != nil {
		return errors.New("数据编码失败")
	}
	err = s.RabbitMqValidate.PublishSimple(string(byteMessage))
	if err != nil {
		return errors.New("数据编码失败")
	}
	a[commodityId].Stock--
	return nil
}
