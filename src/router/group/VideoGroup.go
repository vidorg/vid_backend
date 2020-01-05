package group

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller"
	"github.com/vidorg/vid_backend/src/middleware"
)

func SetupVideoGroup(router *gin.Engine, config *config.ServerConfig) {
	videoCtrl := controller.VideoController(config)

	jwt := middleware.JwtMiddleware(false, config)
	jwtAdmin := middleware.JwtMiddleware(true, config)
	limit := middleware.LimitMiddleware(2 << 20)

	videoGroup := router.Group("/video")
	{
		videoGroup.GET("/", jwtAdmin, videoCtrl.QueryAllVideos)
		videoGroup.GET("/:vid", videoCtrl.QueryVideoByVid)

		videoGroup.POST("/", jwt, limit, videoCtrl.InsertVideo)    // 2M cover
		videoGroup.PUT("/:vid", jwt, limit, videoCtrl.UpdateVideo) // 2M cover
		videoGroup.DELETE("/:vid", jwt, videoCtrl.DeleteVideo)
	}
}
