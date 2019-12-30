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
		authGroup.POST("/login", AuthCtrl.Login)
		authGroup.POST("/register", AuthCtrl.Register)
		authGroup.PUT("/password", jwt, AuthCtrl.ModifyPassword)
		authGroup.GET("/", jwt, AuthCtrl.CurrentUser)
	}
}
