package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type MetaConfig struct {
	RunMode     string `json:"run-mode"      yaml:"run-mode"`
	Port        int32  `json:"port"          yaml:"port"`
	LogPath     string `json:"log-path"      yaml:"log-path"`
	DefPageSize int32  `json:"def-page_size" yaml:"def-page-size"`
	MaxPageSize int32  `json:"max-page_size" yaml:"max-page-size"`
}

type FileConfig struct {
	ImagePath      string `json:"image-path"       yaml:"image-path"`
	ImageMaxSize   int32  `json:"image-max-size"   yaml:"image-max-size"`
	ImageUrlPrefix string `json:"image-url-prefix" yaml:"image-url-prefix"`
}

type MySqlConfig struct {
	Host     string `json:"host"     yaml:"host"`
	Port     int32  `json:"port"     yaml:"port"`
	Name     string `json:"name"     yaml:"name"`
	Charset  string `json:"charset"  yaml:"charset"`
	User     string `json:"user"     yaml:"user"`
	Password string `json:"password" yaml:"password"`
	IsLog    bool   `json:"log"      yaml:"log"`
}

type RedisConfig struct {
	ConnType string `json:"conn-type" yaml:"conn-type"`
	Host     string `json:"host"      yaml:"host"`
	Port     int32  `json:"port"      yaml:"port"`
	Db       int32  `json:"db"        yaml:"db"`
	Password string `json:"password"  yaml:"password"`

	ConnectTimeout int32 `json:"connect-timeout" yaml:"connect-timeout"`
	ReadTimeout    int32 `json:"read-timeout"    yaml:"read-timeout"`
	WriteTimeout   int32 `json:"write-timeout"   yaml:"write-timeout"`
}

type JwtConfig struct {
	Secret   string `json:"secret"    yaml:"secret"`
	Expire   int64  `json:"expire"    yaml:"expire"`
	Issuer   string `json:"issuer"    yaml:"issuer"`
	RedisFmt string `json:"redis-fmt" yaml:"redis-fmt"`
}

type SearchConfig struct {
	DictPath string `json:"dictionary-path" yaml:"dictionary-path"`
}

type ServerConfig struct {
	MetaConfig   *MetaConfig   `yaml:"meta"`
	FileConfig   *FileConfig   `yaml:"file"`
	MySqlConfig  *MySqlConfig  `yaml:"mysql"`
	RedisConfig  *RedisConfig  `yaml:"redis"`
	JwtConfig    *JwtConfig    `yaml:"jwt"`
	SearchConfig *SearchConfig `yaml:"search"`
}

func Load(path string) (*ServerConfig, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &ServerConfig{}
	ext := filepath.Ext(path)
	switch ext {
	case ".yml":
		fallthrough
	case ".yaml":
		err = yaml.Unmarshal(f, config)
	case ".json":
		err = json.Unmarshal(f, config)
	default:
		return nil, fmt.Errorf("Expected a yaml or json file, got a \"%s\" file\n", ext)
	}

	if err != nil {
		return nil, err
	}

	return config, nil
}
