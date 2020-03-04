package conn

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database/helper"
	"log"
	"time"
)

func SetupRedisConn(config *config.RedisConfig) *helper.RedisHelper {
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

	return helper.NewRedisHelper(conn)
}
