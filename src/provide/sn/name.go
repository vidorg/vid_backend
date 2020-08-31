package sn

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
)

const (
	// common
	SConfig xdi.ServiceName = "config" // *config.Config
	SLogger xdi.ServiceName = "logger" // *logrus.Logger

	// databases
	SGorm     xdi.ServiceName = "gorm"            // *gorm.DB
	SRedis    xdi.ServiceName = "redis"           // *redis.Pool
	SEnforcer xdi.ServiceName = "casbin-enforcer" // *casbin.Enforcer

	// services
	SAccountService   xdi.ServiceName = "account-service"   // service.AccountService
	STokenService     xdi.ServiceName = "token-service"     // service.TokenService
	SUserService      xdi.ServiceName = "user-service"      // service.UserService
	SEmailService     xdi.ServiceName = "email-service"     // service.EmailService
	SSubscribeService xdi.ServiceName = "subscribe-service" // service.SubscribeService
	SVideoService     xdi.ServiceName = "video-service"     // service.VideoService
	SJwtService       xdi.ServiceName = "jwt-service"       // service.JwtService
	SCasbinService    xdi.ServiceName = "casbin-service"    // service.CasbinService
)
