package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type MetaConfig struct {
	RunMode     string `yaml:"run-mode"`
	Port        int32  `yaml:"port"`
	LogPath     string `yaml:"log-path"`
	DefPageSize int32  `yaml:"def-page-size"`
	MaxPageSize int32  `yaml:"max-page-size"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int32  `yaml:"port"`
	Name     string `yaml:"name"`
	Charset  string `yaml:"charset"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	IsLog    bool   `yaml:"log"`
}

type RedisConfig struct {
	ConnType string `yaml:"conn-type"`
	Host     string `yaml:"host"`
	Port     int32  `yaml:"port"`
	Db       int32  `yaml:"db"`
	Password string `yaml:"password"`

	ConnectTimeout int32 `yaml:"connect-timeout"`
	ReadTimeout    int32 `yaml:"read-timeout"`
	WriteTimeout   int32 `yaml:"write-timeout"`
}

type JwtConfig struct {
	Secret   string `yaml:"secret"`
	Expire   int64  `yaml:"expire"`
	Issuer   string `yaml:"issuer"`
	RedisFmt string `yaml:"redis-fmt"`
}

type CasbinConfig struct {
	ConfigPath string `yaml:"config-path"`
}

type Config struct {
	Meta   *MetaConfig   `yaml:"meta"`
	MySQL  *MySQLConfig  `yaml:"mysql"`
	Redis  *RedisConfig  `yaml:"redis"`
	Jwt    *JwtConfig    `yaml:"jwt"`
	Casbin *CasbinConfig `yaml:"casbin"`
}

func Load(path string) (*Config, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(f, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
