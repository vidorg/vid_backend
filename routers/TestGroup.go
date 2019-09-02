package routers

import (
	. "vid/controllers"
	"vid/middleware"

	"github.com/gin-gonic/gin"
)

func setupTestGroup(router *gin.Engine) {
	testGroup := router.Group("/test")
	{
		testGroup.GET("/", TestCtrl.Test)
	}

	authTestGroup := router.Group("/authtest")
	{
		authTestGroup.Use(middleware.JWTMiddleware())
		authTestGroup.GET("/", TestCtrl.AuthTest)
	}
}
