package server

import (
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database/dao"
	"github.com/vidorg/vid_backend/src/model/common/profile"
	"github.com/vidorg/vid_backend/src/server/inject"
)

func Inject(config *config.ServerConfig) *inject.Option {
	mapper := profile.CreateMapperProfile(config)

	passDao := dao.PassRepository(config.MySqlConfig)
	tokenDao := dao.TokenRepository(config.RedisConfig, config.JwtConfig.RedisHeader)
	userDao := dao.UserRepository(config.MySqlConfig)
	subDao := dao.SubRepository(config.MySqlConfig)
	videoDao := dao.VideoRepository(config.MySqlConfig)

	return &inject.Option{
		ServerConfig:   config,
		EntityMapper: mapper,

		PassDao:  passDao,
		TokenDao: tokenDao,
		UserDao:  userDao,
		SubDao:   subDao,
		VideoDao: videoDao,
	}
}
