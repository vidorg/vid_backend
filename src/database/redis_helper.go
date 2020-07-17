package database

import (
	"github.com/gomodule/redigo/redis"
)

type RedisHelper struct {
	redis.Conn
}

func NewRedisHelper(conn redis.Conn) *RedisHelper {
	return &RedisHelper{conn}
}

func (db *RedisHelper) DeleteAll(pattern string) bool {
	keys, err := redis.Strings(db.Do("KEYS", pattern))
	if err != nil {
		return false
	}

	cnt := 0
	for _, key := range keys {
		result, err := redis.Int(db.Do("DEL", key))
		if err == nil {
			cnt += result
		}
	}
	return len(keys) == 0 || cnt > 0
}
