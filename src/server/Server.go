package server

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "github.com/vidorg/vid_backend/docs"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/server/router"
	"io"
	"log"
	"net/http"
	"os"
)

func InitServer(config *config.ServerConfig) *http.Server {
	// Gin Server
	engine := gin.Default()
	initLogger()

	gin.SetMode(config.RunMode)
	if config.RunMode == "debug" {
		ginpprof.Wrap(engine)
	}

	// Binding & DI
	SetupDefinedValidation()
	dic := ProvideService(config)

	// Route
	router.SetupCommonRouter(engine)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.SetupV1Router(engine, config, dic)

	// Server
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", config.HTTPConfig.Port),
		Handler: engine,
	}
}

func initLogger() {
	logFile, err := os.Create(fmt.Sprintf("./log/log-%s.log", xstring.CurrentTimeUuid(14)))
	if err != nil {
		log.Fatalln("Failed to create log file:", err)
	}

	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
}
