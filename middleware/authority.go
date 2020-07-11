package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	R "shopping/response"
	"shopping/utils"
)

func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, exists := c.Get("jwtUserInfo")
		if exists {
			info := userInfo.(utils.JwtUserInfo)
			if info.Authority == 2 {
				c.Next()
			} else {
				R.Response(c, http.StatusUnauthorized, "无权限", nil, http.StatusUnauthorized)
				c.Abort()
			}
		} else {
			R.Response(c, http.StatusUnauthorized, "无权限", nil, http.StatusUnauthorized)
			c.Abort()
		}
		return
	}
}
