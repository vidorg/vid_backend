package server

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/ahlib-web/xtime"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/vidorg/vid_backend/docs"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/middleware"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

func init() {
	goapidoc.SetDocument(
		"localhost:3344", "/",
		goapidoc.NewInfo("vid backend", "Vid backend built by golang/gin", "1.2").
			WithTermsOfService("https://github.com/vidorg/vid_backend/issues").
			WithLicense(goapidoc.NewLicense("MIT", "https://github.com/vidorg/vid_backend/blob/master/LICENSE")).
			WithContact(goapidoc.NewContact("vidorg", "https://github.com/vidorg", "")),
	)
	goapidoc.SetTags(
		goapidoc.NewTag("Authorization", "Auth-Controller"),
		goapidoc.NewTag("User", "User-Controller"),
		goapidoc.NewTag("Subscribe", "Subscribe-Controller"),
		goapidoc.NewTag("Video", "Video-Controller"),
		goapidoc.NewTag("Policy", "Policy-Controller"),
		goapidoc.NewTag("Administration", "*-Controller"),
	)
	goapidoc.SetSecurities(
		goapidoc.NewSecurity("Jwt", "header", "Authorization"),
	)
}

type Server struct {
	engine *gin.Engine
	config *config.Config
}

func NewServer() *Server {
	// setting
	cfg := xdi.GetByNameForce(sn.SConfig).(*config.Config)
	gin.SetMode(cfg.Meta.RunMode)
	setupBinding()

	// engine
	engine := gin.New()

	// mw
	engine.Use(middleware.LoggerMiddleware())
	engine.Use(middleware.RecoveryMiddleware())
	engine.Use(middleware.CorsMiddleware())

	// route
	if cfg.Meta.RunMode == "debug" {
		ginpprof.Wrap(engine)
	}
	docs.RegisterSwag()
	swaggerUrl := ginSwagger.URL(fmt.Sprintf("http://localhost:%d/v1/swagger/doc.json", cfg.Meta.Port))
	engine.GET("/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerUrl))
	initRoute(engine)

	// server
	return &Server{engine: engine, config: cfg}
}

func setupBinding() {
	xgin.SetupRegexBinding()

	xgin.SetupDateTimeBinding("date", xtime.RFC3339Date)
	xgin.SetupDateTimeBinding("datetime", xtime.RFC3339DateTime)

	xgin.SetupSpecificRegexpBinding("name", "^[a-zA-Z0-9\u4E00-\u9FBF\u3040-\u30FF\\-_]+$")              // alphabet number character kana - _
	xgin.SetupSpecificRegexpBinding("pwd", "^.+$")                                                       // all
	xgin.SetupSpecificRegexpBinding("phone", "^(13[0-9]|15[012356789]|17[678]|18[0-9]|14[57])[0-9]{8}$") // 11
}

func (s *Server) Serve() error {
	addr := fmt.Sprintf("0.0.0.0:%d", s.config.Meta.Port)
	return s.engine.Run(addr)
}
