package routers

import (
	. "vid/controllers"

	"github.com/gin-gonic/gin"
)

func setupSearchGroup(router *gin.Engine) {
	searchGroup := router.Group("/search")
	{
		searchGroup.GET("/user", SearchCtrl.SearchUser)
		searchGroup.GET("/video", SearchCtrl.SearchVideo)
	}
}
