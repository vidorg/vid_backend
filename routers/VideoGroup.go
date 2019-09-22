package routers

import (
	. "vid/controllers"
	"vid/middleware"

	"github.com/gin-gonic/gin"
)

func setupVideoGroup(router *gin.Engine) {

	jwt := middleware.JWTMiddleware()

	videoGroup := router.Group("/video")
	{
		videoGroup.GET("/all", VideoCtrl.GetAllVideos)
		videoGroup.GET("/uid/:uid", VideoCtrl.GetVideosByUid)
		videoGroup.GET("/vid/:vid", VideoCtrl.GetVideoByVid)

		videoGroup.Use(jwt).POST("/new", VideoCtrl.UploadNewVideo)
		videoGroup.Use(jwt).POST("/update", VideoCtrl.UpdateVideoInfo)
		videoGroup.Use(jwt).DELETE("/delete", VideoCtrl.DeleteVideo)
	}
}
