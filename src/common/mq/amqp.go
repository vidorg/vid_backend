package mq

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/streadway/amqp"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

func NewAmqpConn() (*amqp.Connection, error) {
	cfg := xdi.GetByNameForce(sn.SConfig).(*config.Config)
	acfg := cfg.Amqp

	param := fmt.Sprintf("amqp://%s:%s@%s:%d", acfg.Username, acfg.Password, acfg.Host, acfg.Port)
	conn, err := amqp.Dial(param)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
