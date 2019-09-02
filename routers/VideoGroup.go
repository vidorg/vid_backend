package routers

import (
	"vid/controllers"

	"github.com/gin-gonic/gin"
)

var videoCtrl = new(controllers.VideoCtrl)

func setupVideoGroup(router *gin.Engine) {
	videoGroup := router.Group("/video")
	{
		videoGroup.GET("/all", videoCtrl.GetAllVideos)
		videoGroup.GET("/uid/:uid", videoCtrl.GetVideosByUid)
		videoGroup.GET("/vid/:vid", videoCtrl.GetVideoByVid)
	}
}
