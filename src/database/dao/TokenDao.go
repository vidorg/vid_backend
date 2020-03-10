package dao

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gomodule/redigo/redis"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database/helper"
	"log"
	"strconv"
)

type TokenDao struct {
	Config *config.ServerConfig `di:"~"`
	Conn   *helper.RedisHelper  `di:"~"`
}

func NewTokenDao(dic *xdi.DiContainer) *TokenDao {
	repo := &TokenDao{}
	if !dic.Inject(repo) {
		log.Fatalln("Inject failed")
	}
	return repo
}

func (t *TokenDao) concat(uid string, token string) string {
	return fmt.Sprintf(t.Config.JwtConfig.RedisFmt, uid, token)
}

func (t *TokenDao) Query(token string) bool {
	pattern := t.concat("*", token)
	keys, err := redis.Strings(t.Conn.Do("KEYS", pattern))
	if err != nil {
		return false
	}
	return len(keys) >= 1
}

func (t *TokenDao) Insert(token string, uid int32, ex int64) bool {
	value := t.concat(strconv.Itoa(int(uid)), token)
	_, err := t.Conn.Do("SET", value, uid, "EX", ex)
	return err == nil
}

func (t *TokenDao) Delete(token string) bool {
	pattern := t.concat("*", token)
	return t.Conn.DeleteAll(pattern)
}

func (t *TokenDao) DeleteAll(uid int32) bool {
	pattern := t.concat(strconv.Itoa(int(uid)), "*")
	return t.Conn.DeleteAll(pattern)
}
