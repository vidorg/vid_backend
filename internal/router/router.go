package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/internal/middleware"
	"github.com/vidorg/vid_backend/internal/router/controller"
)

// Init init router
func Init() *gin.Engine {
	r := gin.New()

	//if !viper.GetBool("meta.debug") {
	//	gin.SetMode(gin.ReleaseMode)
	//}

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "msg")
	})
	r.POST("/userLogin", controller.UserLogin)
	r.Use(middleware.Auth()).GET("/auth", func(c *gin.Context) {
		c.String(200, "auth test")
	})
	return r
}
