package routers

import (
	"vid/controllers"

	"github.com/gin-gonic/gin"
)

var userCtrl = new(controllers.UserCtrl)
var subCtrl = new(controllers.SubCtrl)

func setupUserGroup(router *gin.Engine) {
	userGroup := router.Group("/user")
	{
		userGroup.GET("/all", userCtrl.QueryAllUsers)
		userGroup.GET("/one/:uid", userCtrl.QueryUser)
		userGroup.PUT("/insert", userCtrl.InsertUser)
		userGroup.POST("/update", userCtrl.UpdateUser)
		userGroup.DELETE("/delete", userCtrl.DeleteUser)

		userGroup.POST("/sub", subCtrl.SubscribeUser)
		userGroup.POST("/unsub", subCtrl.UnSubscribeUser)
		userGroup.GET("/subscriber/:uid", subCtrl.QuerySubscriberUsers)
		userGroup.GET("/subscribing/:uid", subCtrl.QuerySubscribingUsers)
	}
}
