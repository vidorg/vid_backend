package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/controller"
	"github.com/vidorg/vid_backend/src/middleware"
	"net/http"
)

func initRoute(engine *gin.Engine) {
	// common api
	engine.HandleMethodNotAllowed = true
	engine.NoMethod(func(c *gin.Context) {
		result.Status(http.StatusMethodNotAllowed).JSON(c)
	})
	engine.NoRoute(func(c *gin.Context) {
		result.Status(http.StatusNotFound).SetMessage(fmt.Sprintf("route %s is not found", c.Request.URL.Path)).JSON(c)
	})
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ping": "pong"})
	})

	// service api
	jwtMw := middleware.JwtMiddleware()
	casbinMw := middleware.CasbinMiddleware()
	adminMw := func(c *gin.Context) {
		jwtMw(c)
		casbinMw(c)
		c.Next()
	}

	var (
		authCtrl      = controller.NewAuthController()
		policyCtrl    = controller.NewPolicyController()
		userCtrl      = controller.NewUserController()
		subscribeCtrl = controller.NewSubscribeController()
		videoCtrl     = controller.NewVideoController()
	)
	v1 := engine.Group("/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", authCtrl.Login)
			authGroup.POST("/register", authCtrl.Register)
			authGroup.GET("", adminMw, authCtrl.CurrentUser)
			authGroup.POST("/logout", adminMw, authCtrl.Logout)
			authGroup.PUT("/password", adminMw, authCtrl.UpdatePassword)
		}

		policyGroup := v1.Group("/policy")
		{
			policyGroup.GET("", policyCtrl.Query)
			policyGroup.PUT("/role/:uid", policyCtrl.SetRole)
			policyGroup.POST("/role", policyCtrl.Insert)
			policyGroup.DELETE("/role", policyCtrl.Delete)
		}

		userGroup := v1.Group("/user")
		{
			userGroup.GET("", adminMw, userCtrl.QueryAllUsers)
			userGroup.GET("/:uid", userCtrl.QueryUser)
			userGroup.PUT("", adminMw, userCtrl.UpdateUser(false))
			userGroup.DELETE("", adminMw, userCtrl.DeleteUser(false))

			userGroup.GET("/:uid/subscriber", subscribeCtrl.QuerySubscriberUsers)
			userGroup.GET("/:uid/subscribing", subscribeCtrl.QuerySubscribingUsers)
			userGroup.PUT("/subscribing", adminMw, subscribeCtrl.SubscribeUser)
			userGroup.DELETE("/subscribing", adminMw, subscribeCtrl.UnSubscribeUser)

			userGroup.PUT("/admin/:uid", adminMw, userCtrl.UpdateUser(true))
			userGroup.DELETE("/admin/:uid", adminMw, userCtrl.DeleteUser(true))

			userGroup.GET("/:uid/video", videoCtrl.QueryVideosByUid)
		}

		videoGroup := v1.Group("/video")
		{
			videoGroup.GET("", adminMw, videoCtrl.QueryAllVideos)
			videoGroup.GET("/:vid", videoCtrl.QueryVideoByVid)

			videoGroup.POST("", adminMw, videoCtrl.InsertVideo)
			videoGroup.PUT("/:vid", adminMw, videoCtrl.UpdateVideo)
			videoGroup.DELETE("/:vid", adminMw, videoCtrl.DeleteVideo)
		}
	}
}
