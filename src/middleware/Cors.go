package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
)

func CorsMiddleware(config *config.ServerConfig) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if config.RunMode == "debug" {
				return true
			} else {
				return origin == "http://xxx.yyy.zzz" || origin == "https://xxx.yyy.zzz"
			}
		},
	})
}
