package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/service"
)

func JwtMiddleware(srv *service.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := srv.GetToken(c)
		user, err := srv.JwtCheck(token)
		if err != nil {
			result.Error(err).JSON(c)
			c.Abort()
			return
		}

		c.Set(srv.UserKey, user)
		c.Next()
	}
}
