package router

import (
	"fmt"
	"net/http"
	. "vid/app/controller"
	"vid/app/middleware"
	"vid/app/model/dto"

	"github.com/gin-gonic/gin"
)

func SetupRouters() *gin.Engine {
	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(middleware.CORS(middleware.CORSOptions{}))

	SetupTestGroup(router)
	SetupAuthGroup(router)

	// SetupUserGroup(router)
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

func SetupTestGroup(router *gin.Engine) {
	testGroup := router.Group("/test")
	{
		testGroup.GET("/", TestCtrl.Test)
	}

	authTestGroup := router.Group("/auth-test")
	{
		authTestGroup.Use(middleware.JWTMiddleware())
		authTestGroup.GET("/", TestCtrl.AuthTest)
	}
}
