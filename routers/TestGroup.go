package routers

import (
	"vid/controllers"

	"vid/middleware"

	"github.com/gin-gonic/gin"
)

var testCtrl = new(controllers.TestCtrl)

func setupTestGroup(router *gin.Engine) {
	testGroup := router.Group("/test")
	{
		testGroup.GET("/", testCtrl.Test)
	}

	authTestGroup := router.Group("/authtest")
	{
		authTestGroup.Use(middleware.JWTMiddleware())
		authTestGroup.GET("/", testCtrl.AuthTest)
	}
}
