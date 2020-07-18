package middleware

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"github.com/vidorg/vid_backend/src/service"
)

func CasbinMiddleware() gin.HandlerFunc {
	jwtService := xdi.GetByNameForce(sn.SJwtService).(*service.JwtService)
	casbinService := xdi.GetByNameForce(sn.SCasbinService).(*service.CasbinService)

	return func(c *gin.Context) {
		user := jwtService.GetContextUser(c)
		if user == nil {
			c.Abort()
			result.Error(exception.CheckUserRoleError).JSON(c)
			return
		}
		sub := user.Role
		obj := c.FullPath()
		act := c.Request.Method

		ok, err := casbinService.Enforce(sub, obj, act)
		if err != nil {
			c.Abort()
			result.Error(exception.CheckUserRoleError).JSON(c)
			return
		}
		if !ok {
			c.Abort()
			result.Error(exception.RoleHasNoPermissionError).JSON(c)
			return
		}
		c.Next()
	}
}
