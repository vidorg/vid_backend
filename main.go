package main

import (
	"fmt"
	"net/http"
	"vid/config"
	"vid/database"
	"vid/models"
	"vid/routers"
	"vid/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Config
	cfg := config.LoagServerConfig()
	gin.SetMode(cfg.RunMode)

	// Database
	database.SetupDBConn(cfg)

	// Jwt
	utils.JwtSecret = []byte(cfg.JwtSecret)
	utils.JwtTokenExpire = cfg.JwtTokenExpire

	// range
	models.MinLen_Username = cfg.MinLen_Username
	models.MaxLen_Username = cfg.MaxLen_Username
	models.MinLen_Password = cfg.MinLen_Password
	models.MaxLen_Password = cfg.MaxLen_Password

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
