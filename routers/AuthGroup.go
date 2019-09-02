package routers

import (
	. "vid/controllers"
	"vid/middleware"

	"github.com/gin-gonic/gin"
)

func setupAuthGroup(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", AuthCtrl.Login)
		authGroup.POST("/register", AuthCtrl.Register)
		authGroup.Use(middleware.JWTMiddleware()).POST("/modifypass", AuthCtrl.ModifyPass)
	}
}
