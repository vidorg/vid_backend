package routers

import (
	. "vid/app/controllers"

	"github.com/gin-gonic/gin"
)

func SetupSearchGroup(router *gin.Engine) {
	searchGroup := router.Group("/search")
	{
		searchGroup.GET("/user", SearchCtrl.SearchUser)
		searchGroup.GET("/video", SearchCtrl.SearchVideo)
		searchGroup.GET("/playlist", SearchCtrl.SearchPlaylist)
	}
}
