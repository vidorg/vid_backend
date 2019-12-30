package group

import (
	. "vid/app/controller"
	"vid/app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthGroup(router *gin.Engine) {

	jwt := middleware.JwtMiddleware(false)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", AuthController.Login)
		authGroup.POST("/register", AuthController.Register)
		authGroup.PUT("/password", jwt, AuthController.ModifyPassword)
		authGroup.GET("/", jwt, AuthController.CurrentUser)
	}
}
