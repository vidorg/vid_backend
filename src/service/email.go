package service

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"gopkg.in/gomail.v2"
	"math/rand"
	"time"
)

type EmailService struct {
	config *config.Config
	logger *logrus.Logger
	rpool  *redis.Pool
}

func NewEmailService() *EmailService {
	return &EmailService{
		config: xdi.GetByNameForce(sn.SConfig).(*config.Config),
		logger: xdi.GetByNameForce(sn.SLogger).(*logrus.Logger),
		rpool:  xdi.GetByNameForce(sn.SRedis).(*redis.Pool),
	}
}

func (e *EmailService) GenerateSpec() string {
	rand.Seed(time.Now().UnixNano())
	return xstring.RandLetterNumberString(50)
}

// !! Attention: include url - http://localhost:3344/v1/auth/spec/:spec.
func (e *EmailService) SendTo(to string, spec string) error {
	cfg := e.config.Email

	message := gomail.NewMessage()
	message.SetHeader("From", cfg.Name)
	message.SetHeader("To", to)
	message.SetHeader("Subject", "Active your vid account")

	body := "Please visit the next link to active your account:\n"
	body += "http://localhost:3344/v1/auth/spec/" + spec
	message.SetBody("text/plain", body)

	start := time.Now()
	dialer := gomail.NewDialer(cfg.SmtpHost, int(cfg.SmtpPort), cfg.Username, cfg.Password)
	err := dialer.DialAndSend(message)
	if err != nil {
		return nil
	}

	du := time.Now().Sub(start).String()
	e.logger.WithFields(map[string]interface{}{
		"module":   "gomail",
		"duration": du,
		"to":       to,
	}).Infof(fmt.Sprintf("[Gomail]    OK | %12s | %s", du, to))
	return nil
}

func (e *EmailService) concat(spec string) string {
	return fmt.Sprintf("vid-spec-%s", spec)
}

func (e *EmailService) CheckSpec(spec string) (uint64, bool, error) {
	conn := e.rpool.Get()
	defer conn.Close()

	pattern := e.concat(spec)
	uidString, err := redis.String(conn.Do("GET", pattern))
	if err == redis.ErrNil {
		return 0, false, nil
	} else if err != nil {
		return 0, false, err
	}

	uid, err := xnumber.Atou64(uidString)
	if err != nil {
		return 0, false, err
	}
	return uid, true, nil
}

func (e *EmailService) InsertSpec(uid uint64, spec string) error {
	conn := e.rpool.Get()
	defer conn.Close()

	pattern := e.concat(spec)
	_, err := conn.Do("SET", pattern, uid, "EX", e.config.Email.Expire)
	return err
}

func (e *EmailService) DeleteSpec(spec string) error {
	conn := e.rpool.Get()
	defer conn.Close()

	pattern := e.concat(spec)
	_, err := conn.Do("DEL", pattern)
	return err
}
