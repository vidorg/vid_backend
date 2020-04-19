package server

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/middleware"
	"github.com/vidorg/vid_backend/src/model/profile"
	"github.com/vidorg/vid_backend/src/service"
)

func ProvideServices(config *config.ServerConfig, logger *logrus.Logger) *xdi.DiContainer {
	dic := xdi.NewDiContainer()
	dic.Provide(config)
	dic.Provide(logger)

	entityMappers := profile.CreateEntityMappers(config)
	dic.Provide(entityMappers)
	propertyMappers := profile.CreatePropertyMappers()
	dic.Provide(propertyMappers)

	segSrv := service.NewSegmentService(dic)
	dic.Provide(segSrv)

	gormHelper := database.SetupMySQLConn(config.MySqlConfig, logger)
	dic.Provide(gormHelper) // after config
	redisHelper := database.SetupRedisConn(config.RedisConfig)
	dic.Provide(redisHelper) // after config

	passDao := service.NewAccountService(dic)
	dic.Provide(passDao) // after gorm
	tokenDao := service.NewTokenService(dic)
	dic.Provide(tokenDao) // after redis
	userDao := service.NewUserService(dic)
	dic.Provide(userDao)
	subDao := service.NewSubscribeService(dic)
	dic.Provide(subDao)
	videoDao := service.NewVideoService(dic)
	dic.Provide(videoDao)
	searchDao := service.NewSearchService(dic)
	dic.Provide(searchDao) // after segSrv

	jwtSrv := middleware.NewJwtService(dic)
	dic.Provide(jwtSrv) // after dao

	return dic
}
