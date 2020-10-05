package logger

import (
	"bytes"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

type MQLogHook struct {
	formatter   logrus.Formatter
	amqpChannel *amqp.Channel
}

type MQLogHookConfig struct {
	Formatter logrus.Formatter
}

func NewMQLogHook(config *MQLogHookConfig) (logrus.Hook, error) {
	conn := xdi.GetByNameForce(sn.SAmqp).(*amqp.Connection)
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &MQLogHook{
		formatter:   config.Formatter,
		amqpChannel: ch,
	}, nil
}

func (m *MQLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (m *MQLogHook) Fire(entry *logrus.Entry) error {
	text, err := m.formatter.Format(entry)
	if err != nil {
		return err
	}
	text = bytes.TrimSpace(text)

	pub := amqp.Publishing{
		ContentType: "text/plain",
		Body:        text,
	}
	err = m.amqpChannel.Publish("vid", "vid.logger", false, false, pub)
	if err != nil {
		return err
	}

	return nil
}
