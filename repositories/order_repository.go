package repositories

import (
	"github.com/jinzhu/gorm"
	"shopping/models"
)

type OrderRepository struct {
	Db *gorm.DB
}

type OrderRepositoryImp interface {
	Add(*models.Order) error
	GetSize(int, int) (*[]models.Order, int, error)
}

func (r *OrderRepository) Add(o *models.Order) (err error) {
	err = r.Db.Create(o).Error
	return
}

func (r *OrderRepository) GetSize(start int, size int) (c *[]models.Order, total int, err error) {

	//TODO 连表查询
	c = &[]models.Order{}
	err = r.Db.Limit(size).Offset(start).Preload("Commodity").Preload("User").Find(&c).Error
	err = r.Db.Model(&models.Order{}).Count(&total).Error
	return
}
