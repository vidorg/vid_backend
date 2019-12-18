package main

import (
	"fmt"
	"log"
	"net/http"
	"vid/app/config"
	"vid/app/database"
	"vid/app/router"
	"vid/app/util"

	"github.com/gin-gonic/gin"

	_ "vid/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title vid backend
// @version 1.1
// @description Backend of repo https://github.com/vidorg/vid_vue
// @termsOfService https://github.com/vidorg
// @host localhost:3344
// @basePath /
// @license.name MIT
// @license.url https://github.com/vidorg/vid_backend/blob/master/LICENSE
// @swagger 2.0
func main() {

	// ServerConfig
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	gin.SetMode(cfg.RunMode)

	// Database
	database.SetupDBConn(cfg.DatabaseConfig)

	// Jwt
	util.JwtSecret = []byte(cfg.JwtConfig.Secret)
	util.JwtExpire = cfg.JwtConfig.Expire

	// Router & Middleware
	appRouter := router.SetupRouters()
	appRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// App
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPConfig.Port),
		Handler: appRouter,
	}

	// Run
	log.Printf("Server init on port :%d\n", cfg.HTTPConfig.Port)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
