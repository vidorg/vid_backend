package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"math"
	"time"
)

func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.Request.URL.Path

		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))

		code := c.Writer.Status()
		ip := c.ClientIP()
		length := c.Writer.Size()
		if length < 0 {
			length = 0
		}

		entry := logger.WithFields(logrus.Fields{
			"Method":   method,
			"Path":     path,
			"Code":     code,
			"Length":   length,
			"ClientIP": ip,
			"Latency":  latency,
		})
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("[Gin] %3d | %12s | %15s | %8s | %-5s %s", code, fmt.Sprintf("%dms", latency), ip, fmt.Sprintf("%dB", length), method, path)
			if code >= 500 {
				entry.Error(msg)
			} else if code >= 400 {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		}
	}
}
