package middleware

import (
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwt gin.HandlerFunc, casbin gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwt(c)
		casbin(c)
		c.Next()
	}
}
