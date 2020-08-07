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

func Provide(configPath string) error {
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

	// /src/database/gorm.go
	mysql, err := database.NewMySQLConn()
	if err != nil {
		log.Fatalln("Failed to setup mysql conn:", err)
	}
	xdi.ProvideName(sn.SGorm, mysql)

	// /src/database/gorm.go
	adapter, err := database.NewGormAdapter()
	if err != nil {
		log.Fatalln("Failed to setup mysql adapter:", err)
	}
	xdi.ProvideName(sn.SGormAdapter, adapter)

	// /src/database/redis.go
	redis, err := database.NewRedisConn()
	if err != nil {
		log.Fatalln("Failed to setup redis conn:", err)
	}
	xdi.ProvideName(sn.SRedis, redis)

	// /src/model/profile/profile.go
	profile.BuildEntityMappers()
	profile.BuildPropertyMappers()

	// /src/service/*
	xdi.ProvideName(sn.SAccountService, service.NewAccountService())
	xdi.ProvideName(sn.STokenService, service.NewTokenService())
	xdi.ProvideName(sn.SUserService, service.NewUserService())
	xdi.ProvideName(sn.SSubscribeService, service.NewSubscribeService())
	xdi.ProvideName(sn.SVideoService, service.NewVideoService())
	xdi.ProvideName(sn.SJwtService, service.NewJwtService())
	xdi.ProvideName(sn.SCasbinService, service.NewCasbinService())

	return nil
}
