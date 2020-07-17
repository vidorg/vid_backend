package database

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xredis"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/config"
	"log"
	"time"
)

func SetupRedisConn(config *config.RedisConfig, logger *logrus.Logger) *RedisHelper {
	conn, err := redis.Dial(
		config.ConnType,
		fmt.Sprintf("%s:%d", config.Host, config.Port),
		redis.DialPassword(config.Password),
		redis.DialDatabase(int(config.Db)),
		redis.DialConnectTimeout(time.Duration(config.ConnectTimeout)*time.Millisecond),
		redis.DialReadTimeout(time.Duration(config.ReadTimeout)*time.Millisecond),
		redis.DialWriteTimeout(time.Duration(config.WriteTimeout)*time.Millisecond),
	)
	if err != nil {
		log.Fatalln("Failed to connect redis:", err)
	}

	connLogger := xredis.NewRedisLogrus(conn, logger)
	return NewRedisHelper(connLogger)
}
