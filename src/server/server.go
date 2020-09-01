package server

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/ahlib-web/xvalidator"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/vidorg/vid_backend/docs"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/middleware"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"net/http"
	"regexp"
)

func init() {
	goapidoc.SetDocument(
		"localhost:3344", "/",
		goapidoc.NewInfo("vid backend", "Vid backend built by golang/gin", "1.2").
			TermsOfService("https://github.com/vidorg/vid_backend/issues").
			License(goapidoc.NewLicense("MIT", "https://github.com/vidorg/vid_backend/blob/master/LICENSE")).
			Contact(goapidoc.NewContact("vidorg", "https://github.com/vidorg", "")),
	)

	goapidoc.SetTags(
		goapidoc.NewTag("Authorization", "auth-controller"),
		goapidoc.NewTag("User", "user-controller"),
		goapidoc.NewTag("Subscribe", "subscribe-controller"),
		goapidoc.NewTag("Video", "video-controller"),
		goapidoc.NewTag("Rbac", "rbac-controller"),
		goapidoc.NewTag("Administration", "*-controller"),
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
	cfg := xdi.GetByNameForce(sn.SConfig).(*config.Config)
	gin.SetMode(cfg.Meta.RunMode)
	engine := gin.New()
	setupBinding()

	// mw
	engine.Use(middleware.LoggerMiddleware())
	engine.Use(middleware.RecoveryMiddleware())
	engine.Use(middleware.CorsMiddleware())

	// route
	if cfg.Meta.RunMode == "debug" {
		xgin.PprofWrap(engine)
	}
	docs.RegisterSwag()
	swaggerUrl := ginSwagger.URL(fmt.Sprintf("http://localhost:%d/v1/swagger/doc.json", cfg.Meta.Port))
	engine.GET("/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerUrl))
	engine.GET("/v1/swagger", func(c *gin.Context) { c.Redirect(http.StatusPermanentRedirect, "/v1/swagger/index.html") })
	initRoute(engine)

	// server
	return &Server{engine: engine, config: cfg}
}

func setupBinding() {
	_ = xgin.EnableRegexpBinding()
	_ = xgin.EnableRFC3339DateBinding()
	_ = xgin.EnableRFC3339DateTimeBinding()

	_ = xgin.AddBinding("name", xvalidator.RegexpValidator(regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)))                                         // a-z A-Z 0-9 - _
	_ = xgin.AddBinding("pwd", xvalidator.RegexpValidator(regexp.MustCompile(`^.+$`)))                                                       // all
	_ = xgin.AddBinding("phone", xvalidator.RegexpValidator(regexp.MustCompile(`^(13[0-9]|15[012356789]|17[678]|18[0-9]|14[57])[0-9]{8}$`))) // 11
}

func (s *Server) Serve() error {
	addr := fmt.Sprintf("0.0.0.0:%d", s.config.Meta.Port)
	return s.engine.Run(addr)
}
