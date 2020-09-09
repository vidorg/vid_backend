package middleware

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

func RecoveryMiddleware() gin.HandlerFunc {
	logger := xdi.GetByNameForce(sn.SLogger).(*logrus.Logger)
	skip := 2

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				rid := c.Writer.Header().Get("X-Request-Id")
				r := result.Error(exception.ServerRecoveryError)
				if gin.Mode() == gin.DebugMode {
					others := map[string]interface{}{
						"request_ip": c.ClientIP(),
						"request_id": rid,
					}
					r.Error = xgin.BuildErrorDto(err, c, others, skip, true)
				}

				logger.WithFields(logrus.Fields{
					"module":    "panic",
					"error":     err,
					"requestID": rid,
				}).Errorln(fmt.Sprintf("[Recovery] panic recovered: %v | %s", err, rid))
				r.JSON(c)
			}
		}()
		c.Next()
	}
}
