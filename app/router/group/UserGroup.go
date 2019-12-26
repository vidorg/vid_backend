package group

import (
	. "vid/app/controller"
	"vid/app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserGroup(router *gin.Engine) {

	jwt := middleware.JWTMiddleware(false)
	jwtAdmin := middleware.JWTMiddleware(true)
	limit := middleware.StreamLimitMiddleware(2 << 20)

	userGroup := router.Group("/user")
	{
		userGroup.GET("/", jwtAdmin, UserCtrl.QueryAllUsers)
		userGroup.GET("/:uid", UserCtrl.QueryUser)
		userGroup.PUT("/", jwt, limit, UserCtrl.UpdateUser) // 2M avatar
		userGroup.DELETE("/", jwt, UserCtrl.DeleteUser)

		userGroup.PUT("/subscribing", jwt, SubCtrl.SubscribeUser)
		userGroup.DELETE("/subscribing", jwt, SubCtrl.UnSubscribeUser)

		userGroup.GET("/:uid/subscriber", SubCtrl.QuerySubscriberUsers)
		userGroup.GET("/:uid/subscribing", SubCtrl.QuerySubscribingUsers)
	}
}
