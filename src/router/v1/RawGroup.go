package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller"
	"github.com/vidorg/vid_backend/src/middleware"
)

func SetupRawGroup(api *gin.RouterGroup, config *config.ServerConfig) {
	rawCtrl := controller.RawController(config)

	jwt := middleware.JwtMiddleware(false, config)
	limit := middleware.LimitMiddleware(int64(config.FileConfig.ImageMaxSize << 20)) // MB

	rawGroup := api.Group("/raw")
	{
		rawGroup.POST("/image", jwt, limit, rawCtrl.UploadImage)
		rawGroup.GET("/image/:filename", rawCtrl.RawImage)
	}
}
