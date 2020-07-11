package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const JwtSecret = "kih**&hgyshq##js"

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

type JwtImp interface {
	GenerateToken() (string, error)
	ParseToken(tokens string) (err error)
}

type JwtUserInfo struct {
	Id        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Authority int    `json:"authority"`
}

func (user *JwtUserInfo) GenerateToken() (string, error) {
	claim := jwt.MapClaims{
		"email":     user.Email,
		"id":        user.Id,
		"name":      user.Username,
		"authority": user.Authority,
		"nbf":       time.Now().Unix(),
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Unix() + 3*60*60,
		"iss":       "myxy99.cn",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokens, err := token.SignedString([]byte(JwtSecret))
	return tokens, err
}

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSecret), nil
	}
}

func (user *JwtUserInfo) ParseToken(tokens string) (err error) {
	token, err := jwt.Parse(tokens, secret())
	if err != nil {
		return
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to mapclaim")
		return
	}
	//验证token，如果token被修改过则为false
	if !token.Valid {
		err = errors.New("token is invalid")
		return
	}
	user.Email = claim["email"].(string)
	user.Username = claim["name"].(string)
	user.Authority = claim["authority"].(int)
	user.Id = claim["id"].(uint)
	return err
}
