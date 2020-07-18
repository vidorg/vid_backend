package sn

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
)

const (
	// config and logger
	SConfig xdi.ServiceName = "config"
	SLogger xdi.ServiceName = "logger"

	// database
	SGorm        xdi.ServiceName = "gorm"
	SGormAdapter xdi.ServiceName = "gorm-adapter"
	SRedis       xdi.ServiceName = "redis"

	// service
	SAccountService   xdi.ServiceName = "account-service"
	STokenService     xdi.ServiceName = "token-service"
	SUserService      xdi.ServiceName = "user-service"
	SSubscribeService xdi.ServiceName = "subscribe-service"
	SVideoService     xdi.ServiceName = "video-service"
	SJwtService       xdi.ServiceName = "jwt-service"
	SCasbinService    xdi.ServiceName = "casbin-service"
)
