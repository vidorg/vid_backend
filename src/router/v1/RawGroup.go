package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller"
)

func SetupRawGroup(api *gin.RouterGroup, config *config.ServerConfig) {
	rawCtrl := controller.RawController(config)

	rawGroup := api.Group("/raw")
	{
		rawGroup.GET("/image/:uid/:filename", rawCtrl.RawImage)
	}
}
