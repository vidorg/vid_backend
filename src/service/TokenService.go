package service

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"strconv"
)

type TokenService struct {
	Config *config.ServerConfig  `di:"~"`
	Logger *logrus.Logger        `di:"~"`
	Conn   *database.RedisHelper `di:"~"`
}

func NewTokenService(dic *xdi.DiContainer) *TokenService {
	repo := &TokenService{}
	dic.MustInject(repo)
	return repo
}

func (t *TokenService) concat(uid string, token string) string {
	return fmt.Sprintf(t.Config.JwtConfig.RedisFmt, uid, token)
}

func (t *TokenService) Query(token string) bool {
	pattern := t.concat("*", token)
	keys, err := redis.Strings(t.Conn.Do("KEYS", pattern))
	if err != nil {
		return false
	}
	return len(keys) >= 1
}

func (t *TokenService) Insert(token string, uid int32, ex int64) bool {
	value := t.concat(strconv.Itoa(int(uid)), token)
	_, err := t.Conn.Do("SET", value, uid, "EX", ex)
	return err == nil
}

func (t *TokenService) Delete(token string) bool {
	pattern := t.concat("*", token)
	return t.Conn.DeleteAll(pattern)
}

func (t *TokenService) DeleteAll(uid int32) bool {
	pattern := t.concat(strconv.Itoa(int(uid)), "*")
	return t.Conn.DeleteAll(pattern)
}
