package services

import "shopping/utils"

type SpikeServiceImp interface {
	Shopping(*utils.JwtUserInfo, int) (bool, error)
}

type SpikeService struct {
}

func (s *SpikeService) Shopping(info *utils.JwtUserInfo, commodityId int) (bool, error) {
	return false, nil
}
