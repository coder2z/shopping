package services

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"shopping/models"
	"shopping/utils"
	"strconv"
	"sync"
	"time"
)

type SpikeServiceUri struct {
	Id int `uri:"id" json:"id" binding:"required,numeric"`
}

type SpikeServiceImp interface {
	Shopping(*utils.JwtUserInfo, int, string) error
}

var (
	//锁
	mutex sync.Mutex
)

type SpikeService struct {
	Consistent       utils.ConsistentHashImp
	LocalHost        string
	HostList         []string
	Port             string
	CommodityCache   map[int]*models.Commodity
	RabbitMqValidate *RabbitMQ
}

func (s *SpikeService) Shopping(info *utils.JwtUserInfo, commodityId int, token string) (err error) {
	id := strconv.Itoa(int(info.Id))
	ip, err := s.Consistent.Get(id)
	if err != nil {
		utils.Log.WithFields(log.Fields{"errMsg": err.Error()}).Warningln("hash环获取数据错误")
		return errors.New("选择服务器错误")
	}
	if s.CommodityCache[commodityId].StartTime > time.Now().Unix() {
		return errors.New("商品为开卖！")
	}
	if ip == s.LocalHost {
		mutex.Lock()
		defer mutex.Unlock()
		if s.CommodityCache[commodityId].Stock > 0 {
			s.CommodityCache[commodityId].Stock--
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
		return errors.New("商品已卖完！")
	} else {
		//代理处理
		fmt.Printf(fmt.Sprintf("http://%v:%v/spike/%v", ip, s.Port, commodityId))
		res, _, _ := utils.GetCurl(fmt.Sprintf("http://%v:%v/spike/%v", ip, s.Port, commodityId), token)
		if res.StatusCode == 200 {
			mutex.Lock()
			defer mutex.Unlock()
			s.CommodityCache[commodityId].Stock--
			return nil
		}
		return errors.New("未抢到！")
	}
}
