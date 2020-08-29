package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/controller"
	"github.com/vidorg/vid_backend/src/middleware"
)

func initRoute(engine *gin.Engine) {
	// common api
	engine.HandleMethodNotAllowed = true
	engine.NoMethod(func(c *gin.Context) {
		result.Status(405).JSON(c)
	})
	engine.NoRoute(func(c *gin.Context) {
		result.Status(404).SetMessage(fmt.Sprintf("route %s is not found", c.Request.URL.Path)).JSON(c)
	})
	engine.GET("", func(c *gin.Context) {
		c.JSON(200, &gin.H{"message": "Welcome to vid API."})
	})
	engine.GET("/ping", func(c *gin.Context) {
		result.Ok().SetData(&gin.H{"ping": "pong"}).JSON(c)
	})

	// controller
	v1 := engine.Group("/v1")

	var (
		authCtrl      = controller.NewAuthController()
		userCtrl      = controller.NewUserController()
		subscribeCtrl = controller.NewSubscribeController()
		videoCtrl     = controller.NewVideoController()
		// policyCtrl    = controller.NewPolicyController()
	)

	jwtMw := middleware.JwtMiddleware()
	casbinMw := middleware.CasbinMiddleware()
	authMw := func(c *gin.Context) {
		if !c.IsAborted() {
			jwtMw(c)
			if !c.IsAborted() {
				casbinMw(c)
				if !c.IsAborted() {
					c.Next()
				}
			}
		}
	}

	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/login", authCtrl.Login)
		authGroup.POST("/register", authCtrl.Register)
		authGroup.GET("", authMw, authCtrl.CurrentUser)
		authGroup.POST("/logout", authMw, authCtrl.Logout)
		authGroup.PUT("/password", authMw, authCtrl.UpdatePassword)
	}

	// policyGroup := v1.Group("/policy")
	// {
	// 	policyGroup.Use(authMw)
	// 	policyGroup.GET("", policyCtrl.Query)
	// 	policyGroup.PUT("/:uid/role", policyCtrl.SetRole)
	// 	policyGroup.POST("", policyCtrl.Insert)
	// 	policyGroup.DELETE("", policyCtrl.Delete)
	// }

	userGroup := v1.Group("/user")
	{
		userGroup.GET("", authMw, userCtrl.QueryAllUsers)
		userGroup.GET("/:uid", userCtrl.QueryUser)
		userGroup.PUT("", authMw, userCtrl.UpdateUser(false))
		userGroup.DELETE("", authMw, userCtrl.DeleteUser(false))

		userGroup.GET("/:uid/subscriber", subscribeCtrl.QuerySubscriberUsers)
		userGroup.GET("/:uid/subscribing", subscribeCtrl.QuerySubscribingUsers)
		userGroup.PUT("/subscribing/:uid", authMw, subscribeCtrl.SubscribeUser)
		userGroup.DELETE("/subscribing/:uid", authMw, subscribeCtrl.UnSubscribeUser)

		userGroup.PUT("/admin/:uid", authMw, userCtrl.UpdateUser(true))
		userGroup.DELETE("/admin/:uid", authMw, userCtrl.DeleteUser(true))

		userGroup.GET("/:uid/video", videoCtrl.QueryVideosByUid)
	}

	videoGroup := v1.Group("/video")
	{
		videoGroup.GET("", authMw, videoCtrl.QueryAllVideos)
		videoGroup.GET("/:vid", videoCtrl.QueryVideoByVid)

		videoGroup.POST("", authMw, videoCtrl.InsertVideo)
		videoGroup.PUT("/:vid", authMw, videoCtrl.UpdateVideo)
		videoGroup.DELETE("/:vid", authMw, videoCtrl.DeleteVideo)
	}
}
