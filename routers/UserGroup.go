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
		// public
		userGroup.GET("/all", UserCtrl.QueryAllUsers)
		userGroup.GET("/uid/:uid", UserCtrl.QueryUser)

		userGroup.GET("/subscriber/:uid", SubCtrl.QuerySubscriberUsers)
		userGroup.GET("/subscribing/:uid", SubCtrl.QuerySubscribingUsers)

		// user
		userGroup.Use(jwt).POST("/update", UserCtrl.UpdateUser)
		userGroup.Use(jwt).DELETE("/delete", UserCtrl.DeleteUser)

		userGroup.Use(jwt).POST("/sub", SubCtrl.SubscribeUser)
		userGroup.Use(jwt).POST("/unsub", SubCtrl.UnSubscribeUser)
	}
}
