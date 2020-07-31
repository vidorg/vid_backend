package database

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xredis"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"log"
	"time"
)

func NewRedisConn() (redis.Conn, error) {
	cfg := xdi.GetByNameForce(sn.SConfig).(*config.Config)
	rcfg := cfg.Redis
	logger := xdi.GetByNameForce(sn.SLogger).(*logrus.Logger)

	conn, err := redis.Dial(
		rcfg.ConnType,
		fmt.Sprintf("%s:%d", rcfg.Host, rcfg.Port),
		redis.DialPassword(rcfg.Password),
		redis.DialDatabase(int(rcfg.Db)),
		redis.DialConnectTimeout(time.Duration(rcfg.ConnectTimeout)*time.Millisecond),
		redis.DialReadTimeout(time.Duration(rcfg.ReadTimeout)*time.Millisecond),
		redis.DialWriteTimeout(time.Duration(rcfg.WriteTimeout)*time.Millisecond),
	)
	if err != nil {
		log.Fatalln("Failed to connect redis:", err)
	}

	connLogger := xredis.NewRedisLogrus(conn, logger, cfg.Meta.RunMode == "debug")
	return connLogger, err
}
