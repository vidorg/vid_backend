package routers

import (
	. "vid/controllers"

	"github.com/gin-gonic/gin"
)

func setupVideoGroup(router *gin.Engine) {
	videoGroup := router.Group("/video")
	{
		videoGroup.GET("/all", VideoCtrl.GetAllVideos)
		videoGroup.GET("/uid/:uid", VideoCtrl.GetVideosByUid)
		videoGroup.GET("/vid/:vid", VideoCtrl.GetVideoByVid)
	}
}
