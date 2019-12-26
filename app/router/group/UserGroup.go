package group

import (
	. "vid/app/controller"
	"vid/app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserGroup(router *gin.Engine) {

	jwt := middleware.JWTMiddleware(false)
	jwtAdmin := middleware.JWTMiddleware(true)

	userGroup := router.Group("/user")
	{
		userGroup.GET("/", jwtAdmin, UserCtrl.QueryAllUsers)
		userGroup.GET("/:uid", UserCtrl.QueryUser)
		userGroup.PUT("/", jwt, middleware.StreamLimitMiddleware(2<<20), UserCtrl.UpdateUser) // 2M
		// userGroup.PUT("/", jwt, UserCtrl.UpdateUser) // 2M
		userGroup.DELETE("/", jwt, UserCtrl.DeleteUser)

		userGroup.PUT("/subscribing", jwt, SubCtrl.SubscribeUser)
		userGroup.DELETE("/subscribing", jwt, SubCtrl.UnSubscribeUser)

		userGroup.GET("/:uid/subscriber", SubCtrl.QuerySubscriberUsers)
		userGroup.GET("/:uid/subscribing", SubCtrl.QuerySubscribingUsers)
	}
}
