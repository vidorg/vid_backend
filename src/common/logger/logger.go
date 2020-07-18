package logger

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xlogger"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"time"
)

func Setup() (*logrus.Logger, error) {
	cfg := xdi.GetByNameForce(sn.SConfig).(*config.Config)

	logger := logrus.New()
	logLevel := logrus.WarnLevel
	if cfg.Meta.RunMode == "debug" {
		logLevel = logrus.DebugLevel
	}

	logger.SetLevel(logLevel)
	logger.SetReportCaller(false)
	logger.AddHook(xlogger.NewRotateFileHook(&xlogger.RotateFileConfig{
		MaxSize:    20,
		MaxAge:     30,
		MaxBackups: 15,
		Filename:   cfg.Meta.LogPath,
		Level:      logLevel,
		Formatter:  &logrus.JSONFormatter{TimestampFormat: time.RFC3339},
	}))
	logger.SetFormatter(&xlogger.CustomFormatter{
		ForceColor: true,
	})

	return logger, nil
}
