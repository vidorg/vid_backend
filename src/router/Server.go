package router

import (
	"fmt"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "github.com/vidorg/vid_backend/docs"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/dto/common"
	"github.com/vidorg/vid_backend/src/router/v1"
	"github.com/vidorg/vid_backend/src/util"
	"io"
	"log"
	"net/http"
	"os"
)

func InitServer(config *config.ServerConfig) *http.Server {
	gin.SetMode(config.RunMode)

	// Setup
	initLogger()
	engine := gin.Default()
	commonRouter(engine)

	if config.RunMode == "debug" {
		ginpprof.Wrap(engine)
	}
	SetupDefinedValidation()

	// Route
	v1.SetupRouters(engine, config)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Server
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", config.HTTPConfig.Port),
		Handler: engine,
	}
}

func initLogger() {
	logFile, err := os.Create(fmt.Sprintf("./log/log-%s.log", util.CommonUtil.CurrentTimeString()))
	if err != nil {
		log.Fatalln("Failed to create log file:", err)
	}

	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	// logger := log.New(logFile, "", log.LstdFlags|log.Lshortfile)
}

// @Router				/ping [GET]
// @Summary				Ping
// @Description			Ping
// @Tag					Ping
/* @Success 200			{ "ping": "pong" } */
func commonRouter(router *gin.Engine) {
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
