package routers

import (
	. "vid/controllers"
	"vid/middleware"

	"github.com/gin-gonic/gin"
)

func setupUserGroup(router *gin.Engine) {

	jwt := middleware.JWTMiddleware()

	userGroup := router.Group("/user")
	{
		// Admin
		userGroup.Use(jwt).GET("/all", UserCtrl.QueryAllUsers)

		// Public
		userGroup.GET("/uid/:uid", UserCtrl.QueryUser)

		userGroup.GET("/subscriber/:uid", SubCtrl.QuerySubscriberUsers)
		userGroup.GET("/subscribing/:uid", SubCtrl.QuerySubscribingUsers)

		// Auth
		userGroup.Use(jwt).PUT("/update", UserCtrl.UpdateUser)
		userGroup.Use(jwt).DELETE("/delete", UserCtrl.DeleteUser)

		userGroup.Use(jwt).POST("/sub", SubCtrl.SubscribeUser)
		userGroup.Use(jwt).POST("/unsub", SubCtrl.UnSubscribeUser)
	}
}
