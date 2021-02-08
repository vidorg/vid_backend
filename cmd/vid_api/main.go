package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	conf "github.com/vidorg/vid_backend/internal/config"
	"github.com/vidorg/vid_backend/internal/middleware"
	"github.com/vidorg/vid_backend/internal/router"
	"github.com/vidorg/vid_backend/pkg/jwt"
	"github.com/vidorg/vid_backend/pkg/logger"
	"github.com/vidorg/vid_backend/pkg/orm"
	"github.com/vidorg/vid_backend/pkg/redis"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	flag.Parse()
	if *Help {
		flag.Usage()
	} else {
		run()
	}
}

func run() {

	if err := conf.Load(*Config); err != nil {
		panic(err)
	}

	logger.New(logger.SetEnv("dev"),
		logger.SetDebug(true),
		logger.SetOutput(true),
		logger.SetPath(conf.Config().Meta.LogPath))

	jwt.SetMeta(conf.Config().Jwt.Secret, conf.Config().Jwt.Issuer)

	err := redis.Init(conf.Config().Redis.Addr, conf.Config().Redis.Password, conf.Config().Redis.Db)
	if err != nil {
		logger.Logger().Error("redis initialize err", zap.Error(errors.Wrap(err, "redis initialize err")))
	}

	dbParams := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		conf.Config().MySQL.User, conf.Config().MySQL.Password,
		conf.Config().MySQL.Host, conf.Config().MySQL.Port,
		conf.Config().MySQL.Name, conf.Config().MySQL.Charset,
	)
	if err = orm.Init(mysql.Open(dbParams)); err != nil {
		panic(err)
	}

	engine := router.Init()
	s := &http.Server{
		Addr:           ":" + strconv.Itoa(conf.Config().Meta.Port),
		Handler:        engine,
		MaxHeaderBytes: 1 << 20,
	}

	middleware.Init(engine)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals,
			os.Interrupt,
			os.Kill,
			syscall.SIGQUIT,
			syscall.SIGTERM,
			syscall.SIGINT)
		for {
			sig := <-signals
			switch sig {
			default:
				time.AfterFunc(time.Duration(survivalTimeout), func() {
					logger.Logger().Info(fmt.Sprintf("[%s] shutting down", "vid_api"))
					_ = s.Shutdown(context.Background())
				})
				return
			}
		}
	}()

	if err := s.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}
