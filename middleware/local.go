package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	R "shopping/response"
	"strings"
)

func Local(list []string) gin.HandlerFunc {
	return func(context *gin.Context) {
		ip := strings.Split(context.Request.RemoteAddr, ":")
		for _, v := range list {
			if v == ip[0] {
				context.Next()
				return
			}
		}
		R.Response(context, http.StatusForbidden, "", nil, http.StatusForbidden)
		context.Abort()
		return
	}
}
