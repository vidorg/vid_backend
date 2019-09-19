package routers

import (
	. "vid/controllers"
	"vid/middleware"

	"github.com/gin-gonic/gin"
)

func setupRawGroup(router *gin.Engine) {

	jwt := middleware.JWTMiddleware()

	rawGroup := router.Group("/raw")
	{
		rawGroup.GET("/img/:user/:filename", RawCtrl.RawImg)
		rawGroup.GET("/video/:user/:filename", RawCtrl.RawVideo)
		uploadSubGroup := rawGroup.Group("/upload")
		{
			uploadSubGroup.Use(jwt)
			uploadSubGroup.POST("/img", RawCtrl.UploadImg)
			uploadSubGroup.POST("/video", RawCtrl.UploadVideo)
		}
	}
}
