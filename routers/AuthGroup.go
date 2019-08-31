package routers

import (
	"vid/controllers"

	"github.com/gin-gonic/gin"
)

var authCtrl = new(controllers.AuthCtrl)

func setupAuthGroup(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authCtrl.Login)
		authGroup.POST("/register", authCtrl.Register)
	}
}
