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

// @Router             /ping [GET]
// @Summary            Ping
// @Tag                Ping
// @ResponseDesc 200   "OK"
// @ResponseEx 200     { "ping": "pong" }
func setupCommonRouter(engine *gin.Engine) {
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
}

func setupApiRouter(router *gin.Engine, dic *xdi.DiContainer) {
	container := &struct {
		Config *config.ServerConfig `di:"~"`
		JwtSrv *service.JwtService  `di:"~"`
	}{}
	dic.MustInject(container)

	jwtMw := middleware.JwtMiddleware(container.JwtSrv)
	limit2MMw := middleware.LimitMiddleware(int64(container.Config.FileConfig.ImageMaxSize << 20)) // MB

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
			userGroup.GET("", jwtMw, userCtrl.QueryAllUsers) // admin
			userGroup.GET("/:uid", userCtrl.QueryUser)
			userGroup.PUT("", jwtMw, userCtrl.UpdateUser(false))
			userGroup.DELETE("", jwtMw, userCtrl.DeleteUser(false))

			userGroup.GET("/:uid/video", videoCtrl.QueryVideosByUid)

			userGroup.GET("/:uid/subscriber", subCtrl.QuerySubscriberUsers)
			userGroup.GET("/:uid/subscribing", subCtrl.QuerySubscribingUsers)
			userGroup.PUT("/subscribing", jwtMw, subCtrl.SubscribeUser)
			userGroup.DELETE("/subscribing", jwtMw, subCtrl.UnSubscribeUser)

			userGroup.PUT("/admin/:uid", jwtMw, userCtrl.UpdateUser(true))    // admin
			userGroup.DELETE("/admin/:uid", jwtMw, userCtrl.DeleteUser(true)) // admin
		}

		videoGroup := v1.Group("/video")
		{
			videoGroup.GET("", jwtMw, videoCtrl.QueryAllVideos) // admin
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
