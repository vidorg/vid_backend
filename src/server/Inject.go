package server

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/common/seg"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database/conn"
	"github.com/vidorg/vid_backend/src/database/dao"
	"github.com/vidorg/vid_backend/src/middleware"
	"github.com/vidorg/vid_backend/src/model/profile"
)

func ProvideServices(config *config.ServerConfig, logger *logrus.Logger) *xdi.DiContainer {
	dic := xdi.NewDiContainer()
	dic.Provide(config)
	dic.Provide(logger)

	entityMappers := profile.CreateEntityMappers(config)
	dic.Provide(entityMappers)
	propertyMappers := profile.CreatePropertyMappers()
	dic.Provide(propertyMappers)

	segSrv := seg.NewSegmentService(dic)
	dic.Provide(segSrv)

	gormHelper := conn.SetupMySqlConn(config.MySqlConfig)
	dic.Provide(gormHelper) // after config
	redisHelper := conn.SetupRedisConn(config.RedisConfig)
	dic.Provide(redisHelper) // after config

	passDao := dao.NewPassDao(dic)
	dic.Provide(passDao) // after gorm
	tokenDao := dao.NewTokenDao(dic)
	dic.Provide(tokenDao) // after redis
	userDao := dao.NewUserDao(dic)
	dic.Provide(userDao)
	subDao := dao.NewSubDao(dic)
	dic.Provide(subDao)
	videoDao := dao.NewVideoDao(dic)
	dic.Provide(videoDao)
	searchDao := dao.NewSearchDao(dic)
	dic.Provide(searchDao) // after segSrv

	jwtSrv := middleware.NewJwtService(dic)
	dic.Provide(jwtSrv) // after dao

	return dic
}
