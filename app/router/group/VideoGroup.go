package group

import (
	"github.com/gin-gonic/gin"
	. "vid/app/controller"
	"vid/app/middleware"
)

func SetupVideoGroup(router *gin.Engine) {

	jwt := middleware.JwtMiddleware(false)
	jwtAdmin := middleware.JwtMiddleware(true)
	limit := middleware.StreamLimitMiddleware(2 << 20)

	videoGroup := router.Group("/video")
	{
		videoGroup.GET("/", jwtAdmin, VideoCtrl.QueryAllVideos)
		videoGroup.GET("/:vid", VideoCtrl.QueryVideoByVid)

		videoGroup.POST("/", jwt, limit, VideoCtrl.InsertVideo)    // 2M cover
		videoGroup.PUT("/:vid", jwt, limit, VideoCtrl.UpdateVideo) // 2M cover
		videoGroup.DELETE("/:vid", jwt, VideoCtrl.DeleteVideo)
	}
}
