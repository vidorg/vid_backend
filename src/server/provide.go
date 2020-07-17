package server

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/profile"
	"github.com/vidorg/vid_backend/src/service"
)

func ProvideServices(config *config.Config, logger *logrus.Logger) *xdi.DiContainer {
	dic := xdi.NewDiContainer()

	dic.Provide(config)
	dic.Provide(logger)

	dic.Provide(profile.CreateEntityMappers(config))
	dic.Provide(profile.CreatePropertyMappers())

	dic.Provide(database.SetupMySQLConn(config.MySQL, logger))
	dic.Provide(database.SetupRedisConn(config.Redis, logger))

	dic.Provide(service.NewAccountService(dic))
	dic.Provide(service.NewTokenService(dic))
	dic.Provide(service.NewUserService(dic))
	dic.Provide(service.NewSubscribeService(dic))
	dic.Provide(service.NewVideoService(dic))

	dic.Provide(service.NewSegmentService(dic))
	dic.Provide(service.NewSearchService(dic))
	dic.Provide(service.NewJwtService(dic))
	dic.Provide(service.NewCasbinService(dic))

	return dic
}
