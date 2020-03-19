package router

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller"
	"github.com/vidorg/vid_backend/src/middleware"
	"log"
)

func SetupV1Router(router *gin.Engine, dic *xdi.DiContainer) {
	container := &struct {
		Config *config.ServerConfig   `di:"~"`
		Logger *logrus.Logger         `di:"~"`
		JwtSrv *middleware.JwtService `di:"~"`
	}{}
	if !dic.Inject(container) {
		log.Fatalln("Inject failed")
	}

	corsMw := middleware.CorsMiddleware()
	logger := middleware.LoggerMiddleware(container.Logger)
	jwtMw := container.JwtSrv.JwtMiddleware(false)
	jwtAdminMw := container.JwtSrv.JwtMiddleware(true)
	limit2MMw := middleware.LimitMiddleware(int64(container.Config.FileConfig.ImageMaxSize << 20)) // MB

	router.Use(logger)
	router.Use(gin.Recovery())
	router.Use(corsMw)

	authCtrl := controller.NewAuthController(dic)
	userCtrl := controller.NewUserController(dic)
	subCtrl := controller.NewSubController(dic)
	videoCtrl := controller.NewVideoController(dic)
	rawCtrl := controller.NewRawController(dic)
	searchCtrl := controller.NewSearchController(dic)

	v1 := router.Group("/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", authCtrl.Login)
			authGroup.POST("/register", authCtrl.Register)
			authGroup.GET("", jwtMw, authCtrl.CurrentUser)
			authGroup.POST("/logout", jwtMw, authCtrl.Logout)
			authGroup.PUT("/password", jwtMw, authCtrl.UpdatePassword)
		}

		userGroup := v1.Group("/user")
		{
			userGroup.GET("", jwtAdminMw, userCtrl.QueryAllUsers)
			userGroup.GET("/:uid", userCtrl.QueryUser)
			userGroup.PUT("", jwtMw, userCtrl.UpdateUser(false))
			userGroup.DELETE("", jwtMw, userCtrl.DeleteUser(false))

			userGroup.GET("/:uid/video", videoCtrl.QueryVideosByUid)

			userGroup.GET("/:uid/subscriber", subCtrl.QuerySubscriberUsers)
			userGroup.GET("/:uid/subscribing", subCtrl.QuerySubscribingUsers)
			userGroup.PUT("/subscribing", jwtMw, subCtrl.SubscribeUser)
			userGroup.DELETE("/subscribing", jwtMw, subCtrl.UnSubscribeUser)

			userGroup.PUT("/admin/:uid", jwtAdminMw, userCtrl.UpdateUser(true))
			userGroup.DELETE("/admin/:uid", jwtAdminMw, userCtrl.DeleteUser(true))
		}

		videoGroup := v1.Group("/video")
		{
			videoGroup.GET("", jwtAdminMw, videoCtrl.QueryAllVideos)
			videoGroup.GET("/:vid", videoCtrl.QueryVideoByVid)

			videoGroup.POST("", jwtMw, videoCtrl.InsertVideo)
			videoGroup.PUT("/:vid", jwtMw, videoCtrl.UpdateVideo)
			videoGroup.DELETE("/:vid", jwtMw, videoCtrl.DeleteVideo)
		}

		rawGroup := v1.Group("/raw")
		{
			rawGroup.POST("/image", jwtMw, limit2MMw, rawCtrl.UploadImage)
			rawGroup.GET("/image/:filename", rawCtrl.RawImage)
		}

		searchGroup := v1.Group("/search")
		{
			searchGroup.GET("/user", searchCtrl.SearchUser)
			searchGroup.GET("/video", searchCtrl.SearchVideo)
		}
	}
}
