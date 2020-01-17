package router

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller"
	"github.com/vidorg/vid_backend/src/middleware"
)

func SetupV1Router(router *gin.Engine, config *config.ServerConfig, dic xdi.DiContainer) {

	authCtrl := controller.NewAuthController(dic)
	userCtrl := controller.NewUserController(dic)
	subCtrl := controller.NewSubController(dic)
	videoCtrl := controller.NewVideoController(dic)
	rawCtrl := controller.NewRawController(dic)

	jwt := middleware.JwtMiddleware(false, dic)
	jwtAdmin := middleware.JwtMiddleware(true, dic)
	limit2M := middleware.LimitMiddleware(int64(config.FileConfig.ImageMaxSize << 20)) // MB

	v1 := router.Group("/v1")
	{
		v1.Use(gin.Recovery())
		v1.Use(middleware.CorsMiddleware())

		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", authCtrl.Login)
			authGroup.POST("/register", authCtrl.Register)
			authGroup.GET("/", jwt, authCtrl.CurrentUser)
			authGroup.POST("/logout", jwt, authCtrl.Logout)
			authGroup.PUT("/password", jwt, authCtrl.UpdatePassword)
		}

		userGroup := v1.Group("/user")
		{
			userGroup.GET("/", jwtAdmin, userCtrl.QueryAllUsers)
			userGroup.GET("/:uid", userCtrl.QueryUser)
			userGroup.PUT("/", jwt, userCtrl.UpdateUser(false))
			userGroup.DELETE("/", jwt, userCtrl.DeleteUser(false))

			userGroup.GET("/:uid/video", videoCtrl.QueryVideosByUid)

			userGroup.GET("/:uid/subscriber", subCtrl.QuerySubscriberUsers)
			userGroup.GET("/:uid/subscribing", subCtrl.QuerySubscribingUsers)
			userGroup.PUT("/subscribing", jwt, subCtrl.SubscribeUser)
			userGroup.DELETE("/subscribing", jwt, subCtrl.UnSubscribeUser)

			userGroup.PUT("/admin/:uid", jwtAdmin, userCtrl.UpdateUser(true))
			userGroup.DELETE("/admin/:uid", jwtAdmin, userCtrl.DeleteUser(true))
		}

		videoGroup := v1.Group("/video")
		{
			videoGroup.GET("/", jwtAdmin, videoCtrl.QueryAllVideos)
			videoGroup.GET("/:vid", videoCtrl.QueryVideoByVid)

			videoGroup.POST("/", jwt, videoCtrl.InsertVideo)
			videoGroup.PUT("/:vid", jwt, videoCtrl.UpdateVideo)
			videoGroup.DELETE("/:vid", jwt, videoCtrl.DeleteVideo)
		}

		rawGroup := v1.Group("/raw")
		{
			rawGroup.POST("/image", jwt, limit2M, rawCtrl.UploadImage)
			rawGroup.GET("/image/:filename", rawCtrl.RawImage)
		}
	}
}
