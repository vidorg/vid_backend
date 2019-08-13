package config

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

var Cfg *ini.File

type Config struct {
	RunMode      string
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func LoagServerConfig() Config {
	var err error

	cfg, err := ini.Load("./config/app.ini")
	if err != nil {
		log.Fatal(2, "Fail to parse 'config/app.ini: %v", err)
	}

	server, err := cfg.GetSection("server")
	if err != nil {
		log.Fatal(2, "Fail to parse section 'server': %v", err)
	}

	ret := Config{
		RunMode:      cfg.Section("").Key("RUN_MODE").MustString("debug"),
		HTTPPort:     server.Key("HTTP_PORT").MustInt(),
		ReadTimeout:  time.Duration(server.Key("READ_TIMEOUT").MustInt(60)) * time.Second,
		WriteTimeout: time.Duration(server.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second,
	}

	return ret
}
