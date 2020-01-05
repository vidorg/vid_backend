package group

import (
	. "github.com/vidorg/vid_backend/src/controller"
	"github.com/vidorg/vid_backend/src/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserGroup(router *gin.Engine) {

	jwt := middleware.JwtMiddleware(false)
	jwtAdmin := middleware.JwtMiddleware(true)
	limit := middleware.LimitMiddleware(2 << 20)

	userGroup := router.Group("/user")
	{
		userGroup.GET("/", jwtAdmin, UserController.QueryAllUsers)
		userGroup.GET("/:uid", UserController.QueryUser)
		userGroup.PUT("/", jwt, limit, UserController.UpdateUser) // 2M avatar
		userGroup.DELETE("/", jwt, UserController.DeleteUser)

		userGroup.PUT("/subscribing", jwt, SubController.SubscribeUser)
		userGroup.DELETE("/subscribing", jwt, SubController.UnSubscribeUser)

		userGroup.GET("/:uid/video", VideoController.QueryVideosByUid)

		userGroup.GET("/:uid/subscriber", SubController.QuerySubscriberUsers)
		userGroup.GET("/:uid/subscribing", SubController.QuerySubscribingUsers)
	}
}
