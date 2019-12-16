package router

import (
	. "vid/app/controller"
	"vid/app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserGroup(router *gin.Engine) {

	jwt := middleware.JWTMiddleware()

	userGroup := router.Group("/user")
	{
		// Admin
		userGroup.Use(jwt).GET("/", UserCtrl.QueryAllUsers)

		// Auth
		userGroup.Use(jwt).PUT("/", UserCtrl.UpdateUser)
		userGroup.Use(jwt).DELETE("/", UserCtrl.DeleteUser)

		userGroup.Use(jwt).POST("/sub", SubCtrl.SubscribeUser)
		userGroup.Use(jwt).POST("/unsub", SubCtrl.UnSubscribeUser)

		// Public
		userGroup.GET("/:uid", UserCtrl.QueryUser)
		userGroup.GET("/:uid/subscriber", SubCtrl.QuerySubscriberUsers)
		userGroup.GET("/:uid/subscribing", SubCtrl.QuerySubscribingUsers)
	}
}
