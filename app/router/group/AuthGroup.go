package group

import (
	. "vid/app/controller"
	"vid/app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthGroup(router *gin.Engine) {

	jwt := middleware.JWTMiddleware(false)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", AuthCtrl.Login)
		authGroup.POST("/register", AuthCtrl.Register)
		authGroup.POST("/pass", jwt, AuthCtrl.ModifyPass)
		authGroup.GET("/", jwt, AuthCtrl.CurrentUser)
	}
}
