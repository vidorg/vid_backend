package database

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xredis"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"time"
)

func NewRedisPool() (*redis.Pool, error) {
	cfg := xdi.GetByNameForce(sn.SConfig).(*config.Config)
	rcfg := cfg.Redis
	logger := xdi.GetByNameForce(sn.SLogger).(*logrus.Logger)

	pool := &redis.Pool{
		MaxIdle:         int(rcfg.MaxIdle),
		MaxActive:       int(rcfg.MaxActive),
		MaxConnLifetime: time.Duration(rcfg.MaxLifetime) * time.Second,
		IdleTimeout:     time.Duration(rcfg.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(
				"tcp", fmt.Sprintf("%s:%d", rcfg.Host, rcfg.Port),
				redis.DialPassword(rcfg.Password),
				redis.DialDatabase(int(rcfg.Db)),
				redis.DialConnectTimeout(time.Duration(rcfg.ConnectTimeout)*time.Millisecond),
				redis.DialReadTimeout(time.Duration(rcfg.ReadTimeout)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(rcfg.WriteTimeout)*time.Millisecond),
			)
			if err != nil {
				return nil, err
			}

			conn = xredis.NewLogrusLogger(conn, logger, cfg.Meta.RunMode == "debug").WithSkip(4)
			conn = xredis.NewMutexRedis(conn)
			return conn, nil
		},
	}

	conn := pool.Get()
	defer conn.Close()
	err := conn.Err()
	if err != nil {
		return nil, err
	}

	return pool, nil
}
