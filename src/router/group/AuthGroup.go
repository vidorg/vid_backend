package group

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller"
	"github.com/vidorg/vid_backend/src/middleware"
)

func SetupAuthGroup(router *gin.Engine, config *config.ServerConfig) {
	authCtrl := controller.AuthController(config)

	jwt := middleware.JwtMiddleware(false, config)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authCtrl.Login)
		authGroup.POST("/register", authCtrl.Register)
		authGroup.PUT("/password", jwt, authCtrl.ModifyPassword)
		authGroup.GET("/", jwt, authCtrl.CurrentUser)
	}
}
