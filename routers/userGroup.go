package routers

import (
	"vid/controllers"

	"github.com/gin-gonic/gin"
)

var userCtrl = new(controllers.UserCtrl)

func setupUserGroup(router *gin.Engine) {
	userGroup := router.Group("/user")
	{
		userGroup.GET("/all", userCtrl.QueryAllUsers)
		userGroup.GET("/one/:uid", userCtrl.QueryUser)
		userGroup.PUT("/insert", userCtrl.InsertUser)
		userGroup.POST("/update", userCtrl.UpdateUser)
		userGroup.DELETE("/delete", userCtrl.DeleteUser)
		userGroup.POST("/sub", userCtrl.SubscribeUser)
		userGroup.GET("/subscriber/:uid", userCtrl.QuerySubscriberUsers)
		userGroup.GET("/subscribing/:uid", userCtrl.QuerySubscribingUsers)
	}
}
