package routers

import (
	"vid/controllers"
	"vid/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouters() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(middleware.CORS(middleware.CORSOptions{}))

	setupTestGroup(router)

	// router.NoMethod(controllers.NoMethod)
	router.NoRoute(controllers.NoRoute)

	return router
}
