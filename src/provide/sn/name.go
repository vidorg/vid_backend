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

	// /src/database/redis_conn.go
	SRedis xdi.ServiceName = "redis"
)
