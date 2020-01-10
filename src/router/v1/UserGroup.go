package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller"
	"github.com/vidorg/vid_backend/src/middleware"
)

func SetupUserGroup(api *gin.RouterGroup, config *config.ServerConfig) {
	userCtrl := controller.UserController(config)
	videoCtrl := controller.VideoController(config)
	subCtrl := controller.SubController(config)

	jwt := middleware.JwtMiddleware(false, config)
	jwtAdmin := middleware.JwtMiddleware(true, config)

	userGroup := api.Group("/user")
	{
		userGroup.GET("/", jwtAdmin, userCtrl.QueryAllUsers)
		userGroup.GET("/:uid", userCtrl.QueryUser)
		userGroup.PUT("/", jwt, userCtrl.UpdateUser(false))
		userGroup.DELETE("/", jwt, userCtrl.DeleteUser(false))

		userGroup.GET("/:uid/video", videoCtrl.QueryVideosByUid)

		userGroup.GET("/:uid/subscriber", subCtrl.QuerySubscriberUsers)
		userGroup.GET("/:uid/subscribing", subCtrl.QuerySubscribingUsers)
		userGroup.PUT("/subscribing", jwt, subCtrl.SubscribeUser)
		userGroup.DELETE("/subscribing", jwt, subCtrl.UnSubscribeUser)

		userGroup.PUT("/admin/:uid", jwtAdmin, userCtrl.UpdateUser(true))
		userGroup.DELETE("/admin/:uid", jwtAdmin, userCtrl.DeleteUser(true))
	}
}
