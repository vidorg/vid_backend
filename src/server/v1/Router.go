package v1

import (
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/middleware"
)

func SetupRouters(router *gin.Engine, config *config.ServerConfig, mapper *xmapper.EntitiesMapper) {
	v1 := router.Group("/v1")
	{
		v1.Use(gin.Recovery())
		v1.Use(middleware.CorsMiddleware())

		SetupAuthGroup(v1, config, mapper)
		SetupUserGroup(v1, config, mapper)
		SetupVideoGroup(v1, config, mapper)
		SetupRawGroup(v1, config, mapper)
	}
}
