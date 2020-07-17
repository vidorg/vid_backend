package middleware

import (
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

func RecoveryMiddleware() gin.HandlerFunc {
	lgr := xdi.GetByNameForce(sn.SLogger).(*logrus.Logger)
	skip := 2

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				r := result.Error(exception.ServerRecoveryError)
				if gin.Mode() == gin.DebugMode {
					r.Error = xgin.BuildErrorDto(err, c, skip, true)
				}
				lgr.Errorln("[Recovery] panic recovered:", err)
				r.JSON(c)
			}
		}()
		c.Next()
	}
}
