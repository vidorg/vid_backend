package router

import (
	. "vid/app/controller"
	"vid/app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupVideoGroup(router *gin.Engine) {

	jwt := middleware.JWTMiddleware()

	videoGroup := router.Group("/video")
	{
		// Admin
		videoGroup.Use(jwt).GET("/all", VideoCtrl.GetAllVideos)

		// Public
		videoGroup.GET("/uid/:uid", VideoCtrl.GetVideosByUid)
		videoGroup.GET("/vid/:vid", VideoCtrl.GetVideoByVid)

		// Auth
		videoGroup.Use(jwt).POST("/new", VideoCtrl.UploadNewVideo)
		videoGroup.Use(jwt).PUT("/update", VideoCtrl.UpdateVideoInfo)
		videoGroup.Use(jwt).DELETE("/delete", VideoCtrl.DeleteVideo)
	}
}
