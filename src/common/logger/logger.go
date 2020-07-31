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
	cfg := xdi.GetByNameForce(sn.SConfig).(*config.Config).Meta

	logger := logrus.New()
	logLevel := logrus.WarnLevel
	if cfg.RunMode == "debug" {
		logLevel = logrus.DebugLevel
	}

	logger.SetLevel(logLevel)
	logger.SetReportCaller(false)
	logger.AddHook(xlogger.NewRotateLogHook(&xlogger.RotateLogConfig{
		MaxAge:       15 * 24 * time.Hour,
		RotationTime: 24 * time.Hour,
		LocalTime:    false,
		Filepath:     cfg.LogPath,
		Filename:     cfg.LogName,
		Level:        logLevel,
		Formatter:    &logrus.JSONFormatter{TimestampFormat: time.RFC3339},
	}))

	logger.SetFormatter(&xlogger.CustomFormatter{
		ForceColor: true,
	})

	return logger, nil
}
