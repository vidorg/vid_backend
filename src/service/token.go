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
	rpool *redis.Pool
}

func NewTokenService() *TokenService {
	return &TokenService{
		rpool: xdi.GetByNameForce(sn.SRedis).(*redis.Pool),
	}
}

func (t *TokenService) concat(uid string, token string) string {
	return fmt.Sprintf("vid-token-%s-%s", uid, token)
}

func (t *TokenService) Query(token string) bool {
	conn, err := t.rpool.Dial()
	if err != nil {
		return false
	}
	defer conn.Close()

	pattern := t.concat("*", token)
	keys, err := redis.Strings(conn.Do("KEYS", pattern))
	return err == nil && len(keys) >= 1
}

func (t *TokenService) Insert(token string, uid int32, ex int64) bool {
	conn, err := t.rpool.Dial()
	if err != nil {
		return false
	}
	defer conn.Close()

	pattern := t.concat(xnumber.FormatInt32(uid, 10), token)
	_, err = conn.Do("SET", pattern, uid, "EX", ex)
	return err == nil
}

func (t *TokenService) Delete(token string) bool {
	conn, err := t.rpool.Dial()
	if err != nil {
		return false
	}
	defer conn.Close()

	pattern := t.concat("*", token)
	tot, del, err := xredis.WithConn(conn).DeleteAll(pattern)
	return err == nil && (tot == 0 || del > 0)
}

func (t *TokenService) DeleteAll(uid int32) bool {
	conn, err := t.rpool.Dial()
	if err != nil {
		return false
	}
	defer conn.Close()

	pattern := t.concat(xnumber.FormatInt32(uid, 10), "*")
	tot, del, err := xredis.WithConn(conn).DeleteAll(pattern)
	return err == nil && (tot == 0 || del > 0)
}
