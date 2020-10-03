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
	SCommonService    xdi.ServiceName = "common-service"    // *service.CommonService
	SOrderbyService   xdi.ServiceName = "orderby-service"   // *service.OrderbyService
	SAccountService   xdi.ServiceName = "account-service"   // *service.AccountService
	STokenService     xdi.ServiceName = "token-service"     // *service.TokenService
	SUserService      xdi.ServiceName = "user-service"      // *service.UserService
	SJwtService       xdi.ServiceName = "jwt-service"       // *service.JwtService
	SCasbinService    xdi.ServiceName = "casbin-service"    // *service.CasbinService
	SEmailService     xdi.ServiceName = "email-service"     // *service.EmailService
	SFollowService    xdi.ServiceName = "follow-service"    // *service.FollowService
	SChannelService   xdi.ServiceName = "channel-service"   // *service.ChannelService
	SSubscribeService xdi.ServiceName = "subscribe-service" // *service.SubscribeService
	SVideoService     xdi.ServiceName = "video-service"     // *service.VideoService
	SFavoriteService  xdi.ServiceName = "favorite-service"  // *service.FavoriteService

	// controllers
	SCommonController xdi.ServiceName = "common-controller" // *controller.CommonController
)
