package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type HTTPConfig struct {
	Port int `yaml:"port"`
}

type MySqlConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	Charset  string `yaml:"charset"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	IsLog    bool   `yaml:"log"`
	PageSize int    `yaml:"page-size"`
}

type RedisConfig struct {
	ConnType string `yaml:"conn-type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Db       int    `yaml:"db"`
	Password string `yaml:"password"`

	ConnectTimeout int `yaml:"connect-timeout"`
	ReadTimeout    int `yaml:"read-timeout"`
	WriteTimeout   int `yaml:"write-timeout"`
}

type JwtConfig struct {
	Secret      string `yaml:"secret"`
	Expire      int64  `yaml:"expire"`
	Issuer      string `yaml:"issuer"`
	RedisHeader string `yaml:"redis-header"`
}

type ServerConfig struct {
	RunMode     string       `yaml:"run-mode"`
	HTTPConfig  *HTTPConfig  `yaml:"http"`
	MySqlConfig *MySqlConfig `yaml:"mysql"`
	RedisConfig *RedisConfig `yaml:"redis"`
	JwtConfig   *JwtConfig   `yaml:"jwt"`
}

func Load(configPath string) (*ServerConfig, error) {
	f, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := new(ServerConfig)
	err = yaml.Unmarshal(f, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
