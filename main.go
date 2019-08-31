package main

import (
	"fmt"
	"net/http"
	"vid/config"
	"vid/database"
	"vid/models/head"
	"vid/routers"
	"vid/utils"

	"github.com/gin-gonic/gin"
)

// Setup app config
func setupCfg(cfg config.Config) {
	// Database
	database.SetupDBConn(cfg)

	// Jwt
	utils.JwtSecret = []byte(cfg.JwtSecret)
	utils.JwtTokenExpire = cfg.JwtTokenExpire

	// Format
	head.MinLen_Username = cfg.FormatConfig.MinLen_Username
	head.MaxLen_Username = cfg.FormatConfig.MaxLen_Username
	head.MinLen_Password = cfg.FormatConfig.MinLen_Password
	head.MaxLen_Password = cfg.FormatConfig.MaxLen_Password
}

func main() {
	// Config
	cfg := config.LoagServerConfig()

	gin.SetMode(cfg.RunMode)
	setupCfg(cfg)

	// Router & Middleware
	router := routers.SetupRouters()

	// App
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.HTTPPort),
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		MaxHeaderBytes: 1 << 20,

		Handler: router,
	}

	fmt.Printf("Server init on port :%d\n", cfg.HTTPPort)
	s.ListenAndServe()
}
