package middleware

import (
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	logger := xdi.GetByNameForce(sn.SLogger).(*logrus.Logger)

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		rid := c.Writer.Header().Get("X-Request-ID")
		xgin.WithLogrus(logger, start, c, &xgin.LoggerExtra{
			OtherString: rid,
			OtherFields: map[string]interface{}{"requestID": rid},
		})
	}
}
