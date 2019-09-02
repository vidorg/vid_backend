package routers

import (
	"vid/controllers"

	"github.com/gin-gonic/gin"
)

var videoCtrl = new(controllers.VideoCtrl)

func setupVideoGroup(router *gin.Engine) {
	videoGroup := router.Group("/video")
	{
		videoGroup.GET("/1", videoCtrl.GetVideo)
		videoGroup.GET("/2", videoCtrl.GetUserVideos)
	}
}
