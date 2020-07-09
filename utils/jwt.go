package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const JwtSecret = "jdnsakjbduiiudu"

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

type UserInfo struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
}

func GenerateToken(email string) (string, error) {
	claim := jwt.MapClaims{
		"email": email,
		"nbf":   time.Now().Unix(),
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Unix() + 3*60*60,
		"iss":   "myxy99.cn",
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

func ParseToken(tokens string) (email string, err error) {
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
	email = claim["email"].(string)
	return
}
