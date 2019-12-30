package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type HTTPConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	ConnectionType string `yaml:"type"`
	Host           string `yaml:"host"`
	Name           string `yaml:"name"`
	Charset        string `yaml:"charset"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	PageSize       int    `yaml:"page-size"`
}

type JwtConfig struct {
	Secret string `yaml:"secret"`
	Expire int64  `yaml:"expire"`
}

type ServerConfig struct {
	RunMode        string         `yaml:"run-mode"`
	HTTPConfig     HTTPConfig     `yaml:"http"`
	DatabaseConfig DatabaseConfig `yaml:"database"`
	JwtConfig      JwtConfig      `yaml:"jwt"`
}

func Load() (*ServerConfig, error) {
	f, err := os.Open("./app/config/config.yaml")
	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadAll(f)
	_ = f.Close()
	if err != nil {
		return nil, err
	}

	config := new(ServerConfig)
	err = yaml.Unmarshal(buf, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
