package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller"
	"github.com/vidorg/vid_backend/src/middleware"
)

func SetupVideoGroup(api *gin.RouterGroup, config *config.ServerConfig) {
	videoCtrl := controller.VideoController(config)

	jwt := middleware.JwtMiddleware(false, config)
	jwtAdmin := middleware.JwtMiddleware(true, config)

	videoGroup := api.Group("/video")
	{
		videoGroup.GET("/", jwtAdmin, videoCtrl.QueryAllVideos)
		videoGroup.GET("/:vid", videoCtrl.QueryVideoByVid)

		videoGroup.POST("/", jwt, videoCtrl.InsertVideo)
		videoGroup.PUT("/:vid", jwt, videoCtrl.UpdateVideo)
		videoGroup.DELETE("/:vid", jwt, videoCtrl.DeleteVideo)
	}
}
