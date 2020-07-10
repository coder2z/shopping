package repositories

import (
	"github.com/jinzhu/gorm"
	"shopping/models"
)

type UserRepositoryImp interface {
	GetUserByEmail(email string) (*models.User, error)
	AddUser(user *models.User) error
	UpdateUser(user *models.User) error
	DelUser(id int) error
}

//func NewUserRepository() UserRepositoryImp {
//	return &UserManagerRepository{
//		Db: models.MysqlHandler,
//	}
//}

type UserManagerRepository struct {
	Db *gorm.DB
}

func (u *UserManagerRepository) GetUserByEmail(email string) (user *models.User, err error) {
	user = &models.User{}
	err = u.Db.Where("email=?", email).Find(user).Error
	return
}

func (u *UserManagerRepository) AddUser(user *models.User) (err error) {
	err = u.Db.Create(user).Error
	return
}

func (u *UserManagerRepository) UpdateUser(user *models.User) (err error) {
	err = u.Db.Model(&models.User{}).Update(user).Error
	return
}

func (u *UserManagerRepository) DelUser(id int) (err error) {
	err = u.Db.Where("id=?", id).Delete(&models.User{}).Error
	return
}
