package repositories

import (
	"github.com/jinzhu/gorm"
	"shopping/models"
)

type CommodityRepositoryImp interface {
	//添加商品
	Add(*models.Commodity) error
	//获取商品信息
	GetById(int) (*models.Commodity, error)
	//获取全部商品信息
	GetAll() (*[]models.Commodity, error)
	//获取部分
	//入参 开始位置 多少个 用于分页
	GetSize(int, int) (*[]models.Commodity, int, error)
	//删除商品
	Del(int) error
	//更新商品
	Update(int, *models.Commodity) error
	//商品库存减去1
	UpdateStockMinusOne(int) error
}

type CommodityRepository struct {
	Db *gorm.DB
}

func (r *CommodityRepository) Add(c *models.Commodity) (err error) {
	err = r.Db.Create(c).Error
	return
}

func (r *CommodityRepository) GetById(id int) (c *models.Commodity, err error) {
	c = &models.Commodity{}
	err = r.Db.Where("id=?", id).Find(c).Error
	return
}

func (r *CommodityRepository) GetAll() (c *[]models.Commodity, err error) {
	c = &[]models.Commodity{}
	err = r.Db.Find(&c).Error
	return
}

func (r *CommodityRepository) GetSize(start int, size int) (c *[]models.Commodity, total int, err error) {
	c = &[]models.Commodity{}
	err = r.Db.Limit(size).Offset(start).Find(c).Error
	err = r.Db.Model(&models.Commodity{}).Count(&total).Error
	return
}

func (r *CommodityRepository) Del(id int) (err error) {
	err = r.Db.Where("id=?", id).Delete(&models.Commodity{}).Error
	return
}

func (r *CommodityRepository) Update(id int, c *models.Commodity) (err error) {
	err = r.Db.Model(&models.Commodity{}).Where("id=?", id).Update(c).Error
	return
}

func (r *CommodityRepository) UpdateStockMinusOne(id int) (err error) {
	err = r.Db.Model(&models.Commodity{}).Where("id=?", id).UpdateColumn("stock", gorm.Expr("stock - ?", 1)).Error
	return
}
