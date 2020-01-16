package inject

import (
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database/dao"
)

type Option struct {
	ServerConfig *config.ServerConfig
	EntityMapper *xmapper.EntityMapper

	PassDao  *dao.PassDao
	TokenDao *dao.TokenDao
	UserDao  *dao.UserDao
	SubDao   *dao.SubDao
	VideoDao *dao.VideoDao
}
