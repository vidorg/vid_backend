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
	cfg := xdi.GetByNameForce(sn.SConfig).(*config.Config).Redis
	logger := xdi.GetByNameForce(sn.SLogger).(*logrus.Logger)

	conn, err := redis.Dial(
		cfg.ConnType,
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		redis.DialPassword(cfg.Password),
		redis.DialDatabase(int(cfg.Db)),
		redis.DialConnectTimeout(time.Duration(cfg.ConnectTimeout)*time.Millisecond),
		redis.DialReadTimeout(time.Duration(cfg.ReadTimeout)*time.Millisecond),
		redis.DialWriteTimeout(time.Duration(cfg.WriteTimeout)*time.Millisecond),
	)
	if err != nil {
		log.Fatalln("Failed to connect redis:", err)
	}

	connLogger := xredis.NewRedisLogrus(conn, logger)
	return connLogger, err
}