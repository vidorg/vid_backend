package group

import (
	"github.com/gin-gonic/gin"
	. "github.com/vidorg/vid_backend/src/controller"
)

func SetupRawGroup(router *gin.Engine) {

	rawGroup := router.Group("/raw")
	{
		rawGroup.GET("/image/:uid/:filename", RawController.RawImage)
	}
}
