package main

import (
	"fmt"
	"net/http"
	"vid/config"
	"vid/database"
	"vid/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Config
	cfg := config.LoagServerConfig()
	gin.SetMode(cfg.RunMode)

	// Database
	database.SetupDBConn(cfg)

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

	fmt.Println("Server init on port ", cfg.HTTPPort)
	s.ListenAndServe()
}
