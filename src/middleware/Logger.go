package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

func renderLatency(ns float64) string {
	us := ns / 1e3
	ms := us / 1e3
	s := ms / 1e3
	if s > 1 {
		return fmt.Sprintf("%.4fs", s)
	} else if ms > 1 {
		return fmt.Sprintf("%.4fms", ms)
	} else {
		return fmt.Sprintf("%.4fµs", us)
	}
}

func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.Request.URL.Path

		start := time.Now()
		c.Next()
		stop := time.Now()
		latency := float64(stop.Sub(start).Nanoseconds())

		code := c.Writer.Status()
		ip := c.ClientIP()
		length := c.Writer.Size()
		if length < 0 {
			length = 0
		}

		entry := logger.WithFields(logrus.Fields{
			"Module":   "gin",
			"Method":   method,
			"Path":     path,
			"Code":     code,
			"Length":   length,
			"ClientIP": ip,
			"Latency":  latency,
		})
		if len(c.Errors) == 0 {
			msg := fmt.Sprintf("[Gin] %3d | %12s | %15s | %8s | %-7s %s", code, renderLatency(latency), ip, fmt.Sprintf("%dB", length), method, path)
			if code >= 500 {
				entry.Error(msg)
			} else if code >= 400 {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		} else {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
	}
}
