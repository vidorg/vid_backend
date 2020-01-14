package v1

import (
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller"
	"github.com/vidorg/vid_backend/src/middleware"
)

func SetupAuthGroup(api *gin.RouterGroup, config *config.ServerConfig, mapper *xmapper.EntitiesMapper) {
	authCtrl := controller.AuthController(config, mapper)

	jwt := middleware.JwtMiddleware(false, config)

	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", authCtrl.Login)
		authGroup.POST("/register", authCtrl.Register)
		authGroup.GET("/", jwt, authCtrl.CurrentUser)
		authGroup.POST("/logout", jwt, authCtrl.Logout)
		authGroup.PUT("/password", jwt, authCtrl.UpdatePassword)
	}
}
