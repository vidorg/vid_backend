package main

import (
	"fmt"
	"log"
	"net/http"
	"vid/app/config"
	"vid/app/database"
	"vid/app/routers"
	"vid/app/utils"

	"github.com/gin-gonic/gin"
)

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
	utils.JwtSecret = []byte(cfg.JwtConfig.Secret)
	utils.JwtExpire = cfg.JwtConfig.Expire

	// Router & Middleware
	router := routers.SetupRouters()

	// App
	s := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.HTTPConfig.Port),
		// ReadTimeout:    cfg.HTTPConfig.ReadTimeout,
		// WriteTimeout:   cfg.HTTPConfig.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
		Handler:        router,
	}

	// Run
	log.Printf("Server init on port :%d\n", cfg.HTTPConfig.Port)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
