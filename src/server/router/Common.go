package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/model/common"
	"net/http"
)

// @Router              /ping [GET]
// @Summary             Ping
// @Description         Ping
// @Tag                 Ping
/* @Response 200        {| "Content-Type": "application/json; charset=utf-8" |} { "ping": "pong" } */
func SetupCommonRouter(router *gin.Engine) {
	router.HandleMethodNotAllowed = true

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
