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
		userGroup.PUT("/", jwt, UserCtrl.UpdateUser)
		userGroup.DELETE("/", jwt, UserCtrl.DeleteUser)

		userGroup.POST("/subscribing/:uid", jwt, SubCtrl.SubscribeUser)
		userGroup.DELETE("/subscribing/:uid", jwt, SubCtrl.UnSubscribeUser)

		userGroup.GET("/:uid/subscriber", SubCtrl.QuerySubscriberUsers)
		userGroup.GET("/:uid/subscribing", SubCtrl.QuerySubscribingUsers)
	}
}
