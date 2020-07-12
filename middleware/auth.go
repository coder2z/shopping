package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	R "shopping/response"
	"shopping/utils"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if len(token) < 7 {
			c.Abort()
			R.Response(c, http.StatusUnauthorized, "未登录", nil, http.StatusUnauthorized)
			return
		}
		jwtUserInfo := utils.JwtUserInfo{}
		err := jwtUserInfo.ParseToken(token[7:])
		if err != nil {
			R.Response(c, http.StatusUnauthorized, "未登录", nil, http.StatusUnauthorized)
			c.Abort()
			return
		}
		c.Set("jwtUserInfo", jwtUserInfo)
		c.Next()
		return
	}
}
