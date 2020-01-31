package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/result"
	"net/http"
)

// @Model         Ping
// @Description   Ping
// @Property      ping string true * "pong" pong

// @Router               /ping [GET]
// @Summary              Ping
// @Description          Ping
// @Tag                  Ping
// @ResponseDesc 200     "OK"
// @ResponseHeader 200   { "Content-Type": "application/json; charset=utf-8" }
// @ResponseModel 200    #Ping
// @Response 200         { "ping": "pong" }
func SetupCommonRouter(router *gin.Engine) {
	router.HandleMethodNotAllowed = true

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ping": "pong"})
	})

	router.NoMethod(func(c *gin.Context) {
		result.Result{}.Result(http.StatusMethodNotAllowed).JSON(c)
	})
	router.NoRoute(func(c *gin.Context) {
		result.Result{}.Result(http.StatusNotFound).SetMessage(fmt.Sprintf("route %s is not found", c.Request.URL.Path)).JSON(c)
	})
}
