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
		authGroup.Use(jwt).POST("/pass", AuthCtrl.ModifyPass)
		authGroup.Use(jwt).GET("/", AuthCtrl.CurrentUser)
	}
}
