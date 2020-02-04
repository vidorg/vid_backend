package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3344", "http://127.0.0.1:3344", "https://localhost:3344", "https://127.0.0.1:3344", "*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Authorization", "api_key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowCredentials: true,
	})
}
