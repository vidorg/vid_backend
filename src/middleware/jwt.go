package middleware

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"github.com/vidorg/vid_backend/src/service"
)

func JwtMiddleware() gin.HandlerFunc {
	jwtService := xdi.GetByNameForce(sn.SJwtService).(*service.JwtService)

	return func(c *gin.Context) {
		token := jwtService.GetToken(c)
		user, err := jwtService.JwtCheck(token)
		if err != nil {
			result.Error(err).JSON(c)
			c.Abort()
			return
		}

		c.Set(jwtService.UserKey, user)
		c.Next()
	}
}
