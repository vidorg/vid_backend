package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"vid/app/model/dto/common"
)

type CorsOptions struct {
	Origin string
}

func CorsMiddleware(options CorsOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		// https://github.com/gin-gonic/gin/issues/29#issuecomment-89132826
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
			c.JSON(http.StatusMethodNotAllowed, common.Result{}.Error(http.StatusMethodNotAllowed).SetMessage("method not allowed"))
		} else {
			c.Next()
		}
	}
}
