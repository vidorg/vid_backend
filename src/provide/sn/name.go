package sn

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
)

const (
	SConfig xdi.ServiceName = "config"
	SLogger xdi.ServiceName = "logger"

	SGorm     xdi.ServiceName = "gorm"
	SRedis    xdi.ServiceName = "redis"
	SEnforcer xdi.ServiceName = "casbin-enforcer"

	SAccountService   xdi.ServiceName = "account-service"
	STokenService     xdi.ServiceName = "token-service"
	SUserService      xdi.ServiceName = "user-service"
	SSubscribeService xdi.ServiceName = "subscribe-service"
	SVideoService     xdi.ServiceName = "video-service"
	SJwtService       xdi.ServiceName = "jwt-service"
	SCasbinService    xdi.ServiceName = "casbin-service"
)
