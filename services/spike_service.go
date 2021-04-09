package services

import (
	"context"
	"encoding/json"
	"errors"
	"shopping/constant"
	"shopping/models"
	"shopping/utils"
)

type SpikeServiceUri struct {
	Id  int `uri:"id" json:"id" binding:"required,numeric"`
	UId int `uri:"uid" json:"uid" binding:"required,numeric"`
}

type SpikeServiceImp interface {
	Shopping(*utils.JwtUserInfo, int) error
}

type SpikeService struct {
	RabbitMqValidate *RabbitMQ
}

func (s *SpikeService) Shopping(info *utils.JwtUserInfo, commodityId int) (err error) {
	if !utils.Limit(context.Background(), constant.SpikeKey.Format(commodityId)) {
		return errors.New("商品已经卖完")
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
	return nil
}
