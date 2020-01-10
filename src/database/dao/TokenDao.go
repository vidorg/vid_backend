package dao

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
)

type TokenDao struct {
	config *config.RedisConfig
	conn   redis.Conn
	header string
}

func TokenRepository(config *config.RedisConfig, header string) *TokenDao {
	return &TokenDao{
		config: config,
		conn:   database.SetupRedisConn(config),
		header: header,
	}
}

func (t *TokenDao) catHeader(token string) string {
	return fmt.Sprintf("%s-%s", t.header, token)
}

func (t *TokenDao) Query(token string) bool {
	data := t.catHeader(token)
	n, _ := redis.Int(t.conn.Do("EXISTS", data))
	return n >= 1
}

func (t *TokenDao) Insert(token string, uid int32, ex int64) bool {
	data := t.catHeader(token)
	_, err := t.conn.Do("SET", data, uid, "EX", ex)
	return err == nil
}

func (t *TokenDao) Delete(token string) bool {
	data := t.catHeader(token)
	n, _ := redis.Int(t.conn.Do("DEL", data))
	return n >= 1
}
