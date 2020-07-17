package server

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xdatetime"
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xgin"
	"github.com/gin-gonic/gin"
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"github.com/vidorg/vid_backend/src/config"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"time"
)

func setupBinding() {
	xgin.SetupRegexBinding()

	xgin.SetupDateTimeBinding("date", xdatetime.ISO8601DateFormat)
	xgin.SetupDateTimeBinding("datetime", xdatetime.ISO8601DateTimeFormat)

	xgin.SetupSpecificRegexpBinding("name", "^[a-zA-Z0-9\u4E00-\u9FBF\u3040-\u30FF\\-_]+$")              // alphabet number character kana - _
	xgin.SetupSpecificRegexpBinding("pwd", "^.+$")                                                       // all
	xgin.SetupSpecificRegexpBinding("phone", "^(13[0-9]|15[012356789]|17[678]|18[0-9]|14[57])[0-9]{8}$") // 11
}

func setupLogger(config *config.ServerConfig) *logrus.Logger {
	logger := logrus.New()
	logLevel := logrus.WarnLevel
	if config.MetaConfig.RunMode == "debug" {
		logLevel = logrus.DebugLevel
	}

	// file
	fileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   config.MetaConfig.LogPath,
		MaxSize:    50,
		MaxBackups: 3,
		MaxAge:     30,
		Level:      logLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		},
	})
	if err != nil {
		log.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	// logrus
	logger.SetLevel(logLevel)
	logger.SetReportCaller(true)
	logger.AddHook(fileHook)
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		DisableSorting:  true,
		TimestampFormat: time.RFC3339,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return "", fmt.Sprintf(" %s:%d:", filename, f.Line)
		},
	})

	// writer
	out := io.MultiWriter(ansicolor.NewAnsiColorWriter(os.Stdout))
	log.SetOutput(out)
	gin.DefaultWriter = out
	logger.Out = out

	return logger
}
