package service

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xredis"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gomodule/redigo/redis"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"strconv"
)

type TokenService struct {
	config *config.Config
	conn   redis.Conn
}

func NewTokenService() *TokenService {
	return &TokenService{
		config: xdi.GetByNameForce(sn.SConfig).(*config.Config),
		conn:   xdi.GetByNameForce(sn.SRedis).(redis.Conn),
	}
}

func (t *TokenService) concat(uid string, token string) string {
	// vid-token-%s-%s
	return fmt.Sprintf(t.config.Jwt.RedisFmt, uid, token)
}

func (t *TokenService) Query(token string) bool {
	pattern := t.concat("*", token)
	keys, err := redis.Strings(t.conn.Do("KEYS", pattern))
	return err == nil && len(keys) >= 1
}

func (t *TokenService) Insert(token string, uid int32, ex int64) bool {
	pattern := t.concat(strconv.Itoa(int(uid)), token)
	_, err := t.conn.Do("SET", pattern, uid, "EX", ex)
	return err == nil
}

func (t *TokenService) Delete(token string) bool {
	pattern := t.concat("*", token)
	tot, del, err := xredis.WithConn(t.conn).DeleteAll(pattern)
	return err == nil && (tot == 0 || del > 0)
}

func (t *TokenService) DeleteAll(uid int32) bool {
	pattern := t.concat(strconv.Itoa(int(uid)), "*")
	tot, del, err := xredis.WithConn(t.conn).DeleteAll(pattern)
	return err == nil && (tot == 0 || del > 0)
}
