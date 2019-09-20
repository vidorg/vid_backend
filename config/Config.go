package config

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

// App Global Config
type Config struct {
	RunMode string

	HTTPServer   string
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	DbType     string
	DbHost     string
	DbUser     string
	DbPassword string
	DbName     string

	JwtSecret      string
	JwtTokenExpire int64

	FormatConfig _FormatConfig
}

// App Format Config
type _FormatConfig struct {
	MinLen_Username int
	MaxLen_Username int
	MinLen_Password int
	MaxLen_Password int
}

var AppCfg Config

func LoagServerConfig() {
	var err error

	cfg, err := ini.Load("./config/app.ini")
	if err != nil {
		log.Fatal(2, "Fail to parse 'config/app.ini: %v", err)
	}

	server, err := cfg.GetSection("server")
	if err != nil {
		log.Fatal(2, "Fail to parse section 'server': %v", err)
	}

	database, err := cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to parse section 'database': %v", err)
	}

	jwt, err := cfg.GetSection("jwt")
	if err != nil {
		log.Fatal(2, "Fail to parse section 'jwt': %v", err)
	}

	register, err := cfg.GetSection("register")
	if err != nil {
		log.Fatal(2, "Fail to parse section 'register': %v", err)
	}

	formatConfig := _FormatConfig{
		MinLen_Username: register.Key("MINLEN_USERNAME").MustInt(1),
		MaxLen_Username: register.Key("MAXLEN_USERNAME").MustInt(20),
		MinLen_Password: register.Key("MINLEN_PASSWORD").MustInt(1),
		MaxLen_Password: register.Key("MAXLEN_PASSWORD").MustInt(20),
	}

	AppCfg = Config{
		RunMode: cfg.Section("").Key("RUN_MODE").MustString("debug"),

		HTTPServer:   server.Key("HTTP_SERVER").MustString("127.0.0.1"),
		HTTPPort:     server.Key("HTTP_PORT").MustInt(),
		ReadTimeout:  time.Duration(server.Key("READ_TIMEOUT").MustInt(60)) * time.Second,
		WriteTimeout: time.Duration(server.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second,

		DbType:     database.Key("TYPE").MustString(""),
		DbHost:     database.Key("HOST").MustString(""),
		DbUser:     database.Key("USER").MustString(""),
		DbPassword: database.Key("PASSWORD").MustString(""),
		DbName:     database.Key("NAME").MustString(""),

		JwtSecret:      jwt.Key("JWTSECRET").MustString(""),
		JwtTokenExpire: jwt.Key("TOKEN_EXPIRE").MustInt64(0),

		FormatConfig: formatConfig,
	}
}
