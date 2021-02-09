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
	r := gin.New()

	if !(conf.Config().Meta.RunMode == "debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	r.GET("/api/v1/ping", func(c *gin.Context) {
		c.JSON(200, &serializer.Response{Code: 200, Msg: "pong"})
	})
	r.POST("/api/v1/UserLogin", controller.UserLogin)
	r.POST("/api/v1/UserRegister", controller.UserRegister)
	r.GET("/api/v1/GetCategories", controller.GetCategoryList)
	auth := r.Group("/api/v1/auth").Use(middleware.Auth())
	{
		auth.GET("/UserAuth", controller.AuthUser)
	}
	return r
}
