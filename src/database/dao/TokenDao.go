package dao

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gomodule/redigo/redis"
	"github.com/vidorg/vid_backend/src/config"
)

type TokenDao struct {
	Config *config.ServerConfig `di:"~"`
	Conn   redis.Conn           `di:"~"`

	Header string `di:"-"`
}

func NewTokenDao(dic *xdi.DiContainer) *TokenDao {
	repo := &TokenDao{}
	dic.Inject(repo)
	if xdi.HasNilDi(repo) {
		panic("Has nil di field")
	}

	repo.Header = repo.Config.JwtConfig.RedisHeader
	return repo
}

func (t *TokenDao) catHeader(token string) string {
	return fmt.Sprintf("%s-%s", t.Header, token)
}

func (t *TokenDao) Query(token string) bool {
	data := t.catHeader(token)
	n, _ := redis.Int(t.Conn.Do("EXISTS", data))
	return n >= 1
}

func (t *TokenDao) Insert(token string, uid int32, ex int64) bool {
	data := t.catHeader(token)
	_, err := t.Conn.Do("SET", data, uid, "EX", ex)
	return err == nil
}

func (t *TokenDao) Delete(token string) bool {
	data := t.catHeader(token)
	n, _ := redis.Int(t.Conn.Do("DEL", data))
	return n >= 1
}
