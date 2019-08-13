package main

import (
	"fmt"
	"net/http"
	"vid/config"
	"vid/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Config
	cfg := config.LoagServerConfig()
	gin.SetMode(cfg.RunMode)

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
