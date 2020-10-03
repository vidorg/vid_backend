package provide

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/vidorg/vid_backend/src/common/logger"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"github.com/vidorg/vid_backend/src/service"
	"log"
)

func Provide(configPath string) error {
	// *config.Config
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalln("Failed to load config:", err)
	}
	xdi.ProvideName(sn.SConfig, cfg)

	// *logrus.Logger
	lgr, err := logger.Setup()
	if err != nil {
		log.Fatalln("Failed to setup logger:", err)
	}
	xdi.ProvideName(sn.SLogger, lgr)

	// *gorm.DB
	mysql, err := database.NewMySQLDB()
	if err != nil {
		log.Fatalln("Failed to setup mysql conn:", err)
	}
	xdi.ProvideName(sn.SGorm, mysql)

	// *redis.Pool
	redis, err := database.NewRedisPool()
	if err != nil {
		log.Fatalln("Failed to setup redis conn:", err)
	}
	xdi.ProvideName(sn.SRedis, redis)

	// *casbin.Enforcer
	enforcer, err := database.NewCasbinEnforcer()
	if err != nil {
		log.Fatalln("Failed to setup casbin enforcer:", err)
	}
	xdi.ProvideName(sn.SEnforcer, enforcer)

	// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	// services
	xdi.ProvideName(sn.SCommonService, service.NewCommonService())     // *service.CommonService
	xdi.ProvideName(sn.SOrderbyService, service.NewOrderbyService())   // *service.OrderbyService
	xdi.ProvideName(sn.SAccountService, service.NewAccountService())   // *service.AccountService
	xdi.ProvideName(sn.STokenService, service.NewTokenService())       // *service.TokenService
	xdi.ProvideName(sn.SUserService, service.NewUserService())         // *service.UserService
	xdi.ProvideName(sn.SJwtService, service.NewJwtService())           // *service.JwtService
	xdi.ProvideName(sn.SCasbinService, service.NewCasbinService())     // *service.CasbinService
	xdi.ProvideName(sn.SEmailService, service.NewEmailService())       // *service.EmailService
	xdi.ProvideName(sn.SFollowService, service.NewFollowService())     // *service.FollowService
	xdi.ProvideName(sn.SVideoService, service.NewVideoService())       // *service.VideoService
	xdi.ProvideName(sn.SFavoriteService, service.NewFavoriteService()) // *service.FavoriteService
	xdi.ProvideName(sn.SChannelService, service.NewChannelService())   // *service.ChannelService

	// controllers
	xdi.ProvideName(sn.SCommonController, controller.NewCommonController())

	return nil
}
