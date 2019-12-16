package router

import (
	"fmt"
	"net/http"
	"vid/app/middleware"
	"vid/app/model/dto"
	"vid/app/router/group"

	"github.com/gin-gonic/gin"
)

func SetupRouters() *gin.Engine {
	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(middleware.CORS(middleware.CORSOptions{}))

	group.SetupAuthGroup(router)
	group.SetupUserGroup(router)

	// SetupVideoGroup(router)
	// SetupSearchGroup(router)
	// SetupPlaylistGroup(router)
	// SetupRawGroup(router)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound,
			dto.Result{}.Error(http.StatusNotFound).SetMessage(fmt.Sprintf("Route %s %s is not found", c.Request.Method, c.Request.URL.Path)))
	})

	return router
}