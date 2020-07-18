package server

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller"
	"github.com/vidorg/vid_backend/src/middleware"
	"github.com/vidorg/vid_backend/src/service"
	"net/http"
)

func initRoute(engine *gin.Engine) {
	engine.HandleMethodNotAllowed = true
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ping": "pong"})
	})
	engine.GET("/error", func(c *gin.Context) {
		err := fmt.Errorf("test error in /error")
		result.Error(exception.ServerRecoveryError).SetError(err, c).JSON(c)
	})
	engine.GET("/panic", func(c *gin.Context) {
		panic("test panic in /panic")
	})

	engine.NoMethod(func(c *gin.Context) {
		result.Status(http.StatusMethodNotAllowed).JSON(c)
	})
	engine.NoRoute(func(c *gin.Context) {
		result.Status(http.StatusNotFound).SetMessage(fmt.Sprintf("route %s is not found", c.Request.URL.Path)).JSON(c)
	})

	container := &struct {
		Config        *config.Config         `di:"~"`
		JwtService    *service.JwtService    `di:"~"`
		CasbinService *service.CasbinService `di:"~"`
	}{}
	xdi.MustInject(container)

	jwtMw := middleware.JwtMiddleware(container.JwtService)
	casbinMw := middleware.CasbinMiddleware(container.JwtService, container.CasbinService)
	adminMw := middleware.AuthMiddleware(jwtMw, casbinMw)
	limit2MMw := middleware.LimitMiddleware(int64(container.Config.File.ImageMaxSize << 20)) // MB

	v1 := engine.Group("/v1")
	{
		authCtrl := controller.NewAuthController()
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", authCtrl.Login)
			authGroup.POST("/register", authCtrl.Register)
			authGroup.GET("", adminMw, authCtrl.CurrentUser)
			authGroup.POST("/logout", adminMw, authCtrl.Logout)
			authGroup.PUT("/password", adminMw, authCtrl.UpdatePassword)
		}

		userCtrl := controller.NewUserController()
		subCtrl := controller.NewSubController()
		videoCtrl := controller.NewVideoController()
		userGroup := v1.Group("/user")
		{
			userGroup.GET("", adminMw, userCtrl.QueryAllUsers)
			userGroup.GET("/:uid", userCtrl.QueryUser)
			userGroup.PUT("", adminMw, userCtrl.UpdateUser(false))
			userGroup.DELETE("", adminMw, userCtrl.DeleteUser(false))

			userGroup.GET("/:uid/subscriber", subCtrl.QuerySubscriberUsers)
			userGroup.GET("/:uid/subscribing", subCtrl.QuerySubscribingUsers)
			userGroup.PUT("/subscribing", adminMw, subCtrl.SubscribeUser)
			userGroup.DELETE("/subscribing", adminMw, subCtrl.UnSubscribeUser)

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
