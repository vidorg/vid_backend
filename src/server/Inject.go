package server

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gomodule/redigo/redis"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/database/dao"
	"github.com/vidorg/vid_backend/src/model/common/profile"
)

func ProvideService(config *config.ServerConfig) xdi.DiContainer {
	dic := xdi.NewDiContainer()

	dic.Provide(config)
	mapper := profile.CreateMapperProfile(config)
	dic.Provide(mapper)
	mysqlConn := database.SetupDBConn(config.MySqlConfig)
	dic.Provide(mysqlConn)
	redisConn := database.SetupRedisConn(config.RedisConfig)
	dic.ProvideImpl((*redis.Conn)(nil), redisConn)

	passDao := dao.NewPassDao(dic)
	dic.Provide(passDao)
	tokenDao := dao.NewTokenDao(dic)
	dic.Provide(tokenDao)
	userDao := dao.NewUserDao(dic)
	dic.Provide(userDao)
	subDao := dao.NewSubDao(dic)
	dic.Provide(subDao)
	videoDao := dao.NewVideoDao(dic)
	dic.Provide(videoDao)

	return dic
}
