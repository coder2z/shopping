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
		email, err := utils.ParseToken(token[7:])
		if err != nil {
			R.Response(c, http.StatusUnauthorized, "未登录", nil, http.StatusUnauthorized)
			c.Abort()
			return
		}
		c.Set("email", email)
		c.Next()
		return
	}
}
