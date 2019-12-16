package routers

import (
	. "vid/app/controllers"
	"vid/app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRawGroup(router *gin.Engine) {

	jwt := middleware.JWTMiddleware()

	rawGroup := router.Group("/raw")
	{
		rawGroup.GET("/image/:user/:filename", RawCtrl.RawImage)
		rawGroup.GET("/video/:user/:filename", RawCtrl.RawVideo)
		uploadSubGroup := rawGroup.Group("/upload")
		{
			uploadSubGroup.Use(jwt)
			uploadSubGroup.POST("/image", RawCtrl.UploadImage)
			uploadSubGroup.POST("/video", RawCtrl.UploadVideo)
		}
	}
}
