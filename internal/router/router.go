package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/internal/conf"
	"github.com/vidorg/vid_backend/internal/middleware"
	"github.com/vidorg/vid_backend/internal/router/controller"
	"github.com/vidorg/vid_backend/internal/serializer"
)

// Init init router
func Init() *gin.Engine {
	router := gin.New()

	if !(conf.Config().Meta.RunMode == "debug") {
		gin.SetMode(gin.ReleaseMode)
	}
	r := router.Group("/api/v1")
	{
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, &serializer.Response{Code: 200, Msg: "pong"})
		})
		r.POST("/UserLogin", controller.UserLogin)
		r.POST("/UserRegister", controller.UserRegister)
		r.GET("/GetCategories", controller.GetCategoryList)
		r.GET("/GetVideoList", controller.GetVideoList)
		r.GET("/GetChannelList", controller.GetChannelList)
		auth := r.Group("/auth").Use(middleware.Auth())
		{
			auth.GET("/UserAuth", controller.AuthUser)
		}
	}
	return router
}
