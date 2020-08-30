package middleware

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"github.com/vidorg/vid_backend/src/service"
)

func JwtMiddleware() gin.HandlerFunc {
	jwtService := xdi.GetByNameForce(sn.SJwtService).(service.JwtService)

	return func(c *gin.Context) {
		token := jwtService.GetToken(c)
		user, status, err := jwtService.JwtCheck(token)

		var ex *exception.Error
		switch status {
		case xstatus.JwtSuccess:
			ex = nil
		case xstatus.JwtBlank, xstatus.JwtUserErr, xstatus.JwtNotFound:
			ex = exception.UnAuthorizedError
		case xstatus.JwtInvalid, xstatus.JwtIssuer:
			ex = exception.InvalidTokenError
		case xstatus.JwtExpired:
			ex = exception.TokenExpiredError
		case xstatus.JwtFailed:
			ex = exception.CheckAuthorizeError
		default:
			ex = exception.InvalidTokenError
		}

		if ex != nil {
			result.Error(ex).SetError(err, c).JSON(c)
			c.Abort()
			return
		}

		c.Set(jwtService.UserKey(), user)
	}
}
