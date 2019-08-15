package routers

import (
	"vid/controllers"

	"github.com/gin-gonic/gin"
)

var testCtrl = new(controllers.TestCtrl)

func setupTestGroup(router *gin.Engine) {
	testGroup := router.Group("/test")
	{
		testGroup.GET("/", testCtrl.Test)
	}
}
