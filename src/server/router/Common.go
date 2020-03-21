package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/result"
	"net/http"
)

// @Model         Ping
// @Description   Ping
// @Property      ping string true "pong" (default:pong)

// @Router               /ping [GET]
// @Summary              Ping
// @Tag                  Ping
// @ResponseDesc 200     "OK"
// @ResponseHeader 200   { "Content-Type": "application/json; charset=utf-8" }
// @ResponseModel 200    #Ping
// @ResponseEx 200       { "ping": "pong" }
func SetupCommonRouter(engine *gin.Engine) {
	engine.HandleMethodNotAllowed = true

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ping": "pong"})
	})

	engine.NoMethod(func(c *gin.Context) {
		result.Status(http.StatusMethodNotAllowed).JSON(c)
	})
	engine.NoRoute(func(c *gin.Context) {
		result.Status(http.StatusNotFound).SetMessage(fmt.Sprintf("route %s is not found", c.Request.URL.Path)).JSON(c)
	})
}
