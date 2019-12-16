package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CORSOptions struct {
	Origin string
}

// CORS middleware from https://github.com/gin-gonic/gin/issues/29#issuecomment-89132826
func CORS(options CORSOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1") // allow any origin domain
		if options.Origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", options.Origin)
		}
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method != "GET" && c.Request.Method != "POST" && c.Request.Method != "PUT" && c.Request.Method != "DELETE" {
			c.AbortWithStatus(http.StatusMethodNotAllowed)
		} else {
			c.Next()
		}
	}
}
