package sn

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
)

const (
	// /src/config/config.go
	SConfig xdi.ServiceName = "config"

	// /src/common/logger/logger.go
	SLogger xdi.ServiceName = "logger"

	// /src/database/gorm_conn.go
	SGorm xdi.ServiceName = "gorm"

	// /src/database/gorm_conn.go
	SGormAdapter xdi.ServiceName = "gorm-adapter"

	// /src/database/redis_conn.go
	SRedis xdi.ServiceName = "redis"

	// /src/service/account.go
	SAccountService xdi.ServiceName = "account-service"

	// /src/service/token.go
	STokenService xdi.ServiceName = "token-service"

	// /src/service/user.go
	SUserService xdi.ServiceName = "user-service"

	// /src/service/subscribe.go
	SSubscribeService xdi.ServiceName = "subscribe-service"

	// /src/service/video.go
	SVideoService xdi.ServiceName = "video-service"

	// /src/service/search.go
	SSearchService xdi.ServiceName = "search-service"

	// /src/service/jwt.go
	SJwtService xdi.ServiceName = "jwt-service"

	// /src/service/casbin.go
	SCasbinService xdi.ServiceName = "casbin-service"
)
