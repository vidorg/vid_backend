package server

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gomodule/redigo/redis"
	"github.com/vidorg/vid_backend/src/common/entity"
	"github.com/vidorg/vid_backend/src/common/property"
	"github.com/vidorg/vid_backend/src/common/seg"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/database/dao"
	"github.com/vidorg/vid_backend/src/middleware"
)

func ProvideService(config *config.ServerConfig) *xdi.DiContainer {
	dic := xdi.NewDiContainer()
	dic.Provide(config)

	entityMappers := entity.CreateEntityMappers(config)
	dic.Provide(entityMappers)
	propertyMappers := property.CreatePropertyMappers()
	dic.Provide(propertyMappers)

	segSrv := seg.NewSegmentService(dic)
	dic.Provide(segSrv)

	gormHelper := database.SetupDBConn(config.MySqlConfig)
	dic.Provide(gormHelper) // after config
	redisConn := database.SetupRedisConn(config.RedisConfig)
	dic.ProvideImpl((*redis.Conn)(nil), redisConn) // interface

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
