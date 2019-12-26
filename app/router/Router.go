package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"vid/app/middleware"
	"vid/app/model/dto"
	"vid/app/router/group"
)

func SetupRouters() *gin.Engine {
	router := gin.Default()

	router.HandleMethodNotAllowed = true
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware(middleware.CORSOptions{
		Origin: "",
	}))

	group.SetupAuthGroup(router)
	group.SetupUserGroup(router)
	group.SetupVideoGroup(router)
	// SetupSearchGroup(router)
	// SetupPlaylistGroup(router)
	group.SetupRawGroup(router)

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, dto.Result{}.Error(http.StatusMethodNotAllowed).
			SetMessage("method not allowed"))
	})
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, dto.Result{}.Error(http.StatusNotFound).
			SetMessage(fmt.Sprintf("route %s %s is not found", c.Request.Method, c.Request.URL.Path)))
	})

	return router
}
