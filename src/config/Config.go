package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type HTTPConfig struct {
	Port int32 `yaml:"port"`
}

type FileConfig struct {
	ImagePath      string `yaml:"image-path"`
	ImageMaxSize   int32  `yaml:"image-max-size"`
	ImageUrlPrefix string `yaml:"image-url-prefix"`
}

type MySqlConfig struct {
	Host     string `yaml:"host"`
	Port     int32  `yaml:"port"`
	Name     string `yaml:"name"`
	Charset  string `yaml:"charset"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	IsLog    bool   `yaml:"log"`
	PageSize int32  `yaml:"page-size"`
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

type ServerConfig struct {
	RunMode     string       `yaml:"run-mode"`
	HTTPConfig  *HTTPConfig  `yaml:"http"`
	FileConfig  *FileConfig  `yaml:"file"`
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
