package provide

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/vidorg/vid_backend/src/common/logger"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/profile"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"github.com/vidorg/vid_backend/src/service"
	"log"
)

func Provide(configPath string) {
	// /src/config/config.go
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalln("Failed to load config:", err)
	}
	xdi.ProvideName(sn.SConfig, cfg)

	// /src/common/logger/logger.go
	lgr, err := logger.Setup()
	if err != nil {
		log.Fatalln("Failed to setup logger:", err)
	}
	xdi.ProvideName(sn.SLogger, lgr)

	xdi.ProvideType(profile.CreateEntityMappers(cfg))
	xdi.ProvideType(profile.CreatePropertyMappers())

	xdi.ProvideType(database.SetupMySQLConn(cfg.MySQL, lgr))
	xdi.ProvideType(database.SetupRedisConn(cfg.Redis, lgr))

	xdi.ProvideType(service.NewAccountService())
	xdi.ProvideType(service.NewTokenService())
	xdi.ProvideType(service.NewUserService())
	xdi.ProvideType(service.NewSubscribeService())
	xdi.ProvideType(service.NewVideoService())

	xdi.ProvideType(service.NewSegmentService())
	xdi.ProvideType(service.NewSearchService())
	xdi.ProvideType(service.NewJwtService())
	xdi.ProvideType(service.NewCasbinService())
}
