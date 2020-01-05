package group

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller"
)

func SetupRawGroup(router *gin.Engine, config *config.ServerConfig) {
	rawCtrl := controller.RawController(config)

	rawGroup := router.Group("/raw")
	{
		rawGroup.GET("/image/:uid/:filename", rawCtrl.RawImage)
	}
}
