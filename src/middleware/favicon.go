package middleware

import (
	"github.com/gin-gonic/gin"
)

func FaviconMiddleware() gin.HandlerFunc {
	// https://github.com/thinkerou/favicon/blob/master/favicon.go
	return func(c *gin.Context) {
		if c.Request.RequestURI != "/favicon.ico" {
			return
		}
		if c.Request.Method != "GET" && c.Request.Method != "HEAD" {
			// TODO
			c.Abort()
			return
		}
		c.Header("Content-Type", "image/x-icon")
		// TODO
		// http.ServeContent(c.Writer, c.Request, "", info, reader)
	}
}
