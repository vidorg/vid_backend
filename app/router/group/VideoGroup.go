package group

import (
	"github.com/gin-gonic/gin"
	. "vid/app/controller"
	"vid/app/middleware"
)

func SetupVideoGroup(router *gin.Engine) {

	jwt := middleware.JWTMiddleware(false)
	jwtAdmin := middleware.JWTMiddleware(true)

	videoGroup := router.Group("/video")
	{
		videoGroup.GET("/", jwtAdmin, VideoCtrl.QueryAllVideos)
		videoGroup.GET("/vid/:vid", VideoCtrl.QueryVideoByVid)
		videoGroup.GET("/uid/:uid", VideoCtrl.QueryVideosByUid)

		videoGroup.POST("/", jwt, VideoCtrl.InsertVideo)
		videoGroup.PUT("/:vid", jwt, VideoCtrl.UpdateVideo)
		videoGroup.DELETE("/:vid", jwt, VideoCtrl.DeleteVideo)
	}
}
