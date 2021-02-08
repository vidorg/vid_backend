package conf

import (
	"gopkg.in/yaml.v2"
	"os"
)

// config ...
var config *AppConfig

func Config() *AppConfig {
	if config == nil {
		panic("config not init")
	}
	return config
}

type MetaConfig struct {
	RunMode     string `yaml:"run-mode"`
	Port        int    `yaml:"port"`
	LogPath     string `yaml:"log-path"`
	LogName     string `yaml:"log-name"`
	LogRotate   bool   `yaml:"log-rotate"`
	LogMq       bool   `yaml:"log-mq"`
	DefPageSize int32  `yaml:"def-page-size"`
	MaxPageSize int32  `yaml:"max-page-size"`
}

type MySQLConfig struct {
	Host        string `yaml:"host"`
	Port        int32  `yaml:"port"`
	Name        string `yaml:"name"`
	Charset     string `yaml:"charset"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	MaxIdle     int32  `yaml:"max-idle"`
	MaxActive   int32  `yaml:"max-active"`
	MaxLifetime int32  `yaml:"max-lifetime"`
}

type RedisConfig struct {
	Addr           string `yaml:"addr"`
	Db             int    `yaml:"db"`
	Password       string `yaml:"password"`
	ConnectTimeout int32  `yaml:"connect-timeout"`
	ReadTimeout    int32  `yaml:"read-timeout"`
	WriteTimeout   int32  `yaml:"write-timeout"`
	MaxIdle        int32  `yaml:"max-idle"`
	MaxActive      int32  `yaml:"max-active"`
	MaxLifetime    int32  `yaml:"max-lifetime"`
	IdleTimeout    int32  `yaml:"idle-timeout"`
}

type AmqpConfig struct {
	Host     string `yaml:"host"`
	Port     int32  `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type EmailConfig struct {
	Name     string `yaml:"name"`
	SmtpHost string `yaml:"smtp-host"`
	SmtpPort int32  `yaml:"smtp-port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Expire   int64  `yaml:"expire"`
}

type JwtConfig struct {
	Secret string `yaml:"secret"`
	Expire int64  `yaml:"expire"`
	Issuer string `yaml:"issuer"`
}

type CasbinConfig struct {
	ConfigPath string `yaml:"config-path"`
}

type AppConfig struct {
	Meta   *MetaConfig   `yaml:"meta"`
	MySQL  *MySQLConfig  `yaml:"mysql"`
	Redis  *RedisConfig  `yaml:"redis"`
	Amqp   *AmqpConfig   `yaml:"amqp"`
	Email  *EmailConfig  `yaml:"email"`
	Jwt    *JwtConfig    `yaml:"jwt"`
	Casbin *CasbinConfig `yaml:"casbin"`
}

func Load(path string) error {
	f, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	conf := &AppConfig{}
	err = yaml.Unmarshal(f, conf)
	if err != nil {
		return err
	}

	config = conf
	return nil
}
