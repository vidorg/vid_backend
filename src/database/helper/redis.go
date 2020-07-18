package helper

import (
	"github.com/gomodule/redigo/redis"
)

func RedisDeleteAll(conn redis.Conn, pattern string) (int32, bool) {
	keys, err := redis.Strings(conn.Do("KEYS", pattern))
	if err != nil {
		return 0, false
	}

	cnt := 0
	for _, key := range keys {
		result, err := redis.Int(conn.Do("DEL", key))
		if err == nil {
			cnt += result
		}
	}
	return int32(cnt), len(keys) == 0 || cnt > 0
}
