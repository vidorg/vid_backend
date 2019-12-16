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
	util.JwtSecret = []byte(cfg.JwtConfig.Secret)
	util.JwtExpire = cfg.JwtConfig.Expire

	// Router & Middleware
	router := router.SetupRouters()

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
