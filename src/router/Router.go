package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/middleware"
	"github.com/vidorg/vid_backend/src/model/dto/common"
	"github.com/vidorg/vid_backend/src/router/group"
	"net/http"
)

func SetupRouters(router *gin.Engine) {
	router.HandleMethodNotAllowed = true

	router.Use(gin.Recovery())
	router.Use(middleware.CorsMiddleware())

	group.SetupAuthGroup(router)
	group.SetupUserGroup(router)
	group.SetupVideoGroup(router)
	group.SetupRawGroup(router)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ping": "pong"})
	})
	router.NoMethod(func(c *gin.Context) {
		common.Result{}.Error(http.StatusMethodNotAllowed).SetMessage("method not allowed").JSON(c)
	})
	router.NoRoute(func(c *gin.Context) {
		common.Result{}.Error(http.StatusNotFound).SetMessage(fmt.Sprintf("route %s is not found", c.Request.URL.Path)).JSON(c)
	})
}
