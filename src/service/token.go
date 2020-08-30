package service

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xredis"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/gomodule/redigo/redis"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

type TokenService struct {
	config *config.Config
	rpool  *redis.Pool
}

func NewTokenService() *TokenService {
	return &TokenService{
		config: xdi.GetByNameForce(sn.SConfig).(*config.Config),
		rpool:  xdi.GetByNameForce(sn.SRedis).(*redis.Pool),
	}
}

func (t *TokenService) concat(uid string, token string) string {
	return fmt.Sprintf("vid-token-%s-%s", uid, token)
}

func (t *TokenService) Query(token string) (bool, error) {
	conn, err := t.rpool.Dial()
	if err != nil {
		return false, err
	}
	defer conn.Close()

	pattern := t.concat("*", token)
	keys, err := redis.Strings(conn.Do("KEYS", pattern))

	if err != nil {
		return false, err
	}
	return len(keys) >= 1, nil
}

func (t *TokenService) Insert(token string, uid uint64) error {
	conn, err := t.rpool.Dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	pattern := t.concat(xnumber.U64toa(uid), token)
	_, err = conn.Do("SET", pattern, uid, "EX", t.config.Jwt.Expire)

	return err
}

func (t *TokenService) Delete(token string) error {
	conn, err := t.rpool.Dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	pattern := t.concat("*", token)
	tot, del, err := xredis.WithConn(conn).DeleteAll(pattern)

	if err != nil {
		return err
	} else if tot > 0 && del == 0 {
		return fmt.Errorf("delete token failed")
	}
	return nil
}

func (t *TokenService) DeleteAll(uid uint64) error {
	conn, err := t.rpool.Dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	pattern := t.concat(xnumber.U64toa(uid), "*")
	tot, del, err := xredis.WithConn(conn).DeleteAll(pattern)

	if err != nil {
		return err
	} else if tot > 0 && del == 0 {
		return fmt.Errorf("delete tokens failed")
	}
	return nil
}
