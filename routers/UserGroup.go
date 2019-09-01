package routers

import (
	"vid/controllers"
	"vid/middleware"

	"github.com/gin-gonic/gin"
)

var userCtrl = new(controllers.UserCtrl)
var subCtrl = new(controllers.SubCtrl)

func setupUserGroup(router *gin.Engine) {

	jwt := middleware.JWTMiddleware()

	userGroup := router.Group("/user")
	{
		// public
		userGroup.GET("/all", userCtrl.QueryAllUsers)
		userGroup.GET("/one/:uid", userCtrl.QueryUser)

		userGroup.GET("/subscriber/:uid", subCtrl.QuerySubscriberUsers)
		userGroup.GET("/subscribing/:uid", subCtrl.QuerySubscribingUsers)

		// user
		userGroup.Use(jwt).POST("/update", userCtrl.UpdateUser)
		userGroup.Use(jwt).DELETE("/delete", userCtrl.DeleteUser)

		userGroup.Use(jwt).POST("/sub", subCtrl.SubscribeUser)
		userGroup.Use(jwt).POST("/unsub", subCtrl.UnSubscribeUser)
	}
}
