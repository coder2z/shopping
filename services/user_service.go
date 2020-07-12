package services

import (
	"errors"
	"shopping/models"
	"shopping/repositories"
	"shopping/utils"
)

type UserLoginService struct {
	Email    string `form:"email" json:"email" binding:"required,min=5,max=30"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=40"`
}

type UserRegisterService struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Tel      string `json:"tel" form:"tel" binding:"required,len=11"`
	UserName string `json:"user_name" form:"userName" binding:"required,min=5,max=30"`
	PassWord string `json:"pass_word" form:"password" binding:"required,min=5,max=30"`
}

type UserServiceImp interface {
	Login(*UserLoginService) (token string, err error)
	Register(*UserRegisterService) error
	Info(email string) (*models.User, error)
}

type UserService struct {
	UserRepository repositories.UserRepositoryImp `inject:""`
}

//func NewUserServices(repository repositories.UserRepositoryImp) UserServiceImp {
//	return &UserService{repository}
//}

func (u *UserService) Login(server *UserLoginService) (token string, err error) {
	userInfo, err := u.UserRepository.GetUserByEmail(server.Email)
	if err != nil {
		return "", errors.New("用户名或者密码错误")
	}
	if ok := userInfo.CheckPassword(server.Password); !ok {
		return "", errors.New("用户名或者密码错误")
	}
	jwtUserInfo := utils.JwtUserInfo{Email: userInfo.Email, Id: int(userInfo.ID), Username: userInfo.UserName, Authority: userInfo.Authority}
	token, err = jwtUserInfo.GenerateToken()
	if err != nil {
		return "", errors.New("token生成失败")
	}
	return
}

func (u *UserService) Register(server *UserRegisterService) (err error) {
	user := &models.User{
		Email:    server.Email,
		Tel:      server.Tel,
		PassWord: server.PassWord,
		UserName: server.UserName,
	}
	if err = user.SetPassword(user.PassWord); err != nil {
		return errors.New("加密失败")
	}
	if err = u.UserRepository.AddUser(user); err != nil {
		return errors.New("注册失败,用户存在")
	}
	return
}

func (u *UserService) Info(email string) (user *models.User, err error) {
	user, err = u.UserRepository.GetUserByEmail(email)
	return
}
