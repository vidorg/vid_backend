package service

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xredis"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/gomodule/redigo/redis"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

type TokenService struct {
	conn redis.Conn
}

func NewTokenService() *TokenService {
	return &TokenService{
		conn: xdi.GetByNameForce(sn.SRedis).(redis.Conn),
	}
}

func (t *TokenService) concat(uid string, token string) string {
	return fmt.Sprintf("vid-token-%s-%s", uid, token)
}

func (t *TokenService) Query(token string) bool {
	pattern := t.concat("*", token)
	keys, err := redis.Strings(t.conn.Do("KEYS", pattern))
	return err == nil && len(keys) >= 1
}

func (t *TokenService) Insert(token string, uid int32, ex int64) bool {
	pattern := t.concat(xnumber.FormatInt32(uid, 10), token)
	_, err := t.conn.Do("SET", pattern, uid, "EX", ex)
	return err == nil
}

func (t *TokenService) Delete(token string) bool {
	pattern := t.concat("*", token)
	tot, del, err := xredis.WithConn(t.conn).DeleteAll(pattern)
	return err == nil && (tot == 0 || del > 0)
}

func (t *TokenService) DeleteAll(uid int32) bool {
	pattern := t.concat(xnumber.FormatInt32(uid, 10), "*")
	tot, del, err := xredis.WithConn(t.conn).DeleteAll(pattern)
	return err == nil && (tot == 0 || del > 0)
}
