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
		userGroup.Use(jwtAdmin).GET("/", UserCtrl.QueryAllUsers)
		userGroup.Use(jwt).PUT("/", UserCtrl.UpdateUser)
		userGroup.Use(jwt).DELETE("/", UserCtrl.DeleteUser)
		userGroup.Use(jwt).POST("/sub", SubCtrl.SubscribeUser)
		userGroup.Use(jwt).POST("/unsub", SubCtrl.UnSubscribeUser)

		userGroup.GET("/:uid", UserCtrl.QueryUser)
		userGroup.GET("/:uid/subscriber", SubCtrl.QuerySubscriberUsers)
		userGroup.GET("/:uid/subscribing", SubCtrl.QuerySubscribingUsers)
	}
}
