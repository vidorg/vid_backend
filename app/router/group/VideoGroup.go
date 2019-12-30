package group

import (
	"github.com/gin-gonic/gin"
	. "vid/app/controller"
	"vid/app/middleware"
)

func SetupVideoGroup(router *gin.Engine) {

	jwt := middleware.JwtMiddleware(false)
	jwtAdmin := middleware.JwtMiddleware(true)
	limit := middleware.LimitMiddleware(2 << 20)

	videoGroup := router.Group("/video")
	{
		videoGroup.GET("/", jwtAdmin, VideoController.QueryAllVideos)
		videoGroup.GET("/:vid", VideoController.QueryVideoByVid)

		videoGroup.POST("/", jwt, limit, VideoController.InsertVideo)    // 2M cover
		videoGroup.PUT("/:vid", jwt, limit, VideoController.UpdateVideo) // 2M cover
		videoGroup.DELETE("/:vid", jwt, VideoController.DeleteVideo)
	}
}
