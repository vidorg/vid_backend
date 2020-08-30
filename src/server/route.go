package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/controller"
	"github.com/vidorg/vid_backend/src/middleware"
)

func j(fn func(c *gin.Context) *result.Result) func(c *gin.Context) {
	return result.J(fn)
}

func initRoute(engine *gin.Engine) {
	// common api
	engine.HandleMethodNotAllowed = true
	engine.NoRoute(func(c *gin.Context) {
		result.Status(404).SetMessage(fmt.Sprintf("route %s is not found", c.Request.URL.Path)).JSON(c)
	})
	engine.NoMethod(func(c *gin.Context) {
		result.Status(405).SetMessage("method not allowed").JSON(c)
	})
	engine.GET("", func(c *gin.Context) {
		c.JSON(200, &gin.H{"message": "Welcome to vid API."})
	})
	engine.GET("ping", func(c *gin.Context) {
		c.JSON(200, &gin.H{"ping": "pong"})
	})

	// controller
	v1 := engine.Group("v1")

	var (
		authCtrl      = controller.NewAuthController()
		userCtrl      = controller.NewUserController()
		subscribeCtrl = controller.NewSubscribeController()
		videoCtrl     = controller.NewVideoController()
		rbacCtrl      = controller.NewRbacController()
	)

	jwtMw := middleware.JwtMiddleware()
	casbinMw := middleware.CasbinMiddleware()
	authMw := func(c *gin.Context) {
		if !c.IsAborted() {
			jwtMw(c)
			if !c.IsAborted() {
				casbinMw(c)
			}
		}
	}

	authGroup := v1.Group("auth")
	{
		authGroup.POST("register", j(authCtrl.Register))
		authGroup.POST("login", j(authCtrl.Login))
		authGroup.GET("user", authMw, j(authCtrl.CurrentUser))
		authGroup.DELETE("logout", authMw, j(authCtrl.Logout))
		authGroup.PUT("password", authMw, j(authCtrl.UpdatePassword))
	}

	userGroup := v1.Group("user")
	{
		userGroup.GET("", authMw, j(userCtrl.QueryAll))
		userGroup.GET(":uid", j(userCtrl.QueryByUid))
		userGroup.PUT("", authMw, j(userCtrl.Update))
		userGroup.DELETE("", authMw, j(userCtrl.Delete))
		userGroup.GET(":uid/subscriber", j(subscribeCtrl.QuerySubscribers))
		userGroup.GET(":uid/subscribing", j(subscribeCtrl.QuerySubscribings))
		userGroup.PUT("subscribing/:uid", authMw, j(subscribeCtrl.SubscribeUser))
		userGroup.DELETE("subscribing/:uid", authMw, j(subscribeCtrl.UnSubscribeUser))

		userGroup.GET(":uid/video", j(videoCtrl.QueryVideosByUid))
	}

	videoGroup := v1.Group("video")
	{
		videoGroup.GET("", authMw, j(videoCtrl.QueryAllVideos))
		videoGroup.GET(":vid", j(videoCtrl.QueryVideoByVid))
		videoGroup.POST("", authMw, j(videoCtrl.InsertVideo))
		videoGroup.PUT(":vid", authMw, j(videoCtrl.UpdateVideo))
		videoGroup.DELETE(":vid", authMw, j(videoCtrl.DeleteVideo))
	}

	rbacGroup := v1.Group("rbac")
	{
		rbacGroup.Use(jwtMw, casbinMw)
		rbacGroup.GET("rule", j(rbacCtrl.GetRules))
		rbacGroup.PUT("user/{uid}", j(rbacCtrl.ChangeRole))
		rbacGroup.POST("subject", j(rbacCtrl.InsertSubject))
		rbacGroup.POST("policy", j(rbacCtrl.InsertPolicy))
		rbacGroup.DELETE("subject", j(rbacCtrl.DeleteSubject))
		rbacGroup.DELETE("policy", j(rbacCtrl.DeletePolicy))
	}
}
