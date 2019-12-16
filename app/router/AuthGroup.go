package router

import (
	. "vid/app/controller"
	"vid/app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthGroup(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", AuthCtrl.Login)
		authGroup.POST("/register", AuthCtrl.Register)
		authGroup.Use(middleware.JWTMiddleware()).POST("/pass", AuthCtrl.ModifyPass)
	}
}
