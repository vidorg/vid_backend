package routers

import (
	"fmt"
	"net/http"
	"vid/middleware"
	. "vid/models"

	"github.com/gin-gonic/gin"
)

func SetupRouters() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(middleware.CORS(middleware.CORSOptions{}))

	setupTestGroup(router)
	setupUserGroup(router)

	// router.NoMethod(controllers.NoMethod)
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, Message{
			Message: fmt.Sprintf("Route %s %s is not found.", c.Request.Method, c.Request.URL.Path),
		})
	})

	return router
}
