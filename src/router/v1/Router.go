package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/middleware"
)

func SetupRouters(router *gin.Engine, config *config.ServerConfig) {
	v1 := router.Group("/v1")
	{
		v1.Use(gin.Recovery())
		v1.Use(middleware.CorsMiddleware())

		SetupAuthGroup(v1, config)
		SetupUserGroup(v1, config)
		SetupVideoGroup(v1, config)
		SetupRawGroup(v1, config)
	}
}
