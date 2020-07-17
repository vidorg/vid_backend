package server

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xtime"
	"github.com/Aoi-hosizora/ahlib-web/xvalidator"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "github.com/vidorg/vid_backend/docs"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/middleware"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

type Server struct {
	engine *gin.Engine
	config *config.Config
}

func NewServer() *Server {
	cfg := xdi.GetByNameForce(sn.SConfig).(*config.Config)
	gin.SetMode(cfg.Meta.RunMode)
	setupBinding()

	engine := gin.New()

	// mw
	engine.Use(middleware.LoggerMiddleware())
	engine.Use(middleware.RecoveryMiddleware())
	engine.Use(middleware.CorsMiddleware())

	// route
	if cfg.Meta.RunMode == "debug" {
		ginpprof.Wrap(engine)
	}
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	initRoute(engine)

	// server
	return &Server{engine: engine, config: cfg}
}

func (s *Server) Serve() error {
	addr := fmt.Sprintf("0.0.0.0:%d", s.config.Meta.Port)
	return s.engine.Run(addr)
}

func setupBinding() {
	xvalidator.SetupRegexBinding()

	xvalidator.SetupDateTimeBinding("date", xtime.RFC3339Date)
	xvalidator.SetupDateTimeBinding("datetime", xtime.RFC3339DateTime)

	xvalidator.SetupSpecificRegexpBinding("name", "^[a-zA-Z0-9\u4E00-\u9FBF\u3040-\u30FF\\-_]+$")              // alphabet number character kana - _
	xvalidator.SetupSpecificRegexpBinding("pwd", "^.+$")                                                       // all
	xvalidator.SetupSpecificRegexpBinding("phone", "^(13[0-9]|15[012356789]|17[678]|18[0-9]|14[57])[0-9]{8}$") // 11
}
