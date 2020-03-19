package server

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "github.com/vidorg/vid_backend/docs"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/server/router"
	"log"
	"net/http"
)

type Server struct {
	Server *http.Server
	Config *config.ServerConfig
	Dic    *xdi.DiContainer
}

func NewServer(config *config.ServerConfig) *Server {
	// Gin Server & Binding
	engine := gin.New()

	gin.SetMode(config.RunMode)
	if config.RunMode == "debug" {
		ginpprof.Wrap(engine)
	}
	SetupBinding()

	// DI
	logger := SetupLogger(config)
	dic := ProvideServices(config, logger)

	// Route
	router.SetupV1Router(engine, dic)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.SetupCommonRouter(engine)

	// Server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.MetaConfig.Port),
		Handler: engine,
	}
	return &Server{
		Server: server,
		Config: config,
		Dic:    dic,
	}
}

func (s *Server) Serve() {
	fmt.Println()
	log.Printf("Server init on port %s\n\n", s.Server.Addr)
	if err := s.Server.ListenAndServe(); err != nil {
		log.Fatalln("Failed to listen and serve:", err)
	}
}
