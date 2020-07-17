package server

import (
	"fmt"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "github.com/vidorg/vid_backend/docs"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/middleware"
	"log"
	"net/http"
)

type Server struct {
	Engine *gin.Engine
	Config *config.Config
}

func NewServer(config *config.Config) *Server {
	engine := gin.New()
	gin.SetMode(config.Meta.RunMode)

	// setup
	setupBinding()
	logger := setupLogger(config)
	dic := ProvideServices(config, logger)

	// mw
	engine.Use(middleware.LoggerMiddleware(logger))
	engine.Use(middleware.RecoveryMiddleware(config, logger))
	engine.Use(middleware.CorsMiddleware(config))

	// route
	if config.Meta.RunMode == "debug" {
		ginpprof.Wrap(engine)
	}
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	setupApiRouter(engine, dic)
	setupCommonRouter(engine)

	// server
	return &Server{
		Engine: engine,
		Config: config,
	}
}

func (s *Server) Serve() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Config.Meta.Port),
		Handler: s.Engine,
	}

	fmt.Println()
	log.Printf("Server is listening on port %s\n\n", server.Addr)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to listen port and serve: %v\n", err)
	}
}
