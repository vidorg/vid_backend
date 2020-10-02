package xgorm

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"strings"
	"time"
)

func NewGormLogger(logger *logrus.Logger, config logger.Config) logger.Interface {
	return &gormLogger{Logger: logger, Config: config}
}

type gormLogger struct {
	*logrus.Logger
	logger.Config
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l gormLogger) trim(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\n")
	return s
}

func (l gormLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	duration := elapsed.String()
	sql, rows := fc()
	source := utils.FileWithLineNum()

	field := l.WithFields(logrus.Fields{
		"module":   "gorm",
		"type":     "sql",
		"source":   source,
		"duration": duration,
		"sql":      sql,
		"rows":     rows,
	})

	switch {
	case err != nil && l.LogLevel >= logger.Error:
		if sql == "" {
			sql = "?"
		}
		field.Errorf(fmt.Sprintf("[Gorm] %s | %s | %s", err, sql, source))
	case l.LogLevel >= logger.Info:
		if strings.HasPrefix(sql, "SELECT DATABASE()") ||
			strings.HasPrefix(sql, "SELECT count(*) FROM information_schema.tables WHERE table_schema") ||
			strings.HasPrefix(sql, "SELECT count(*) FROM information_schema.statistics WHERE table_schema") {
			break
		}
		field.Infof(fmt.Sprintf("[Gorm] #: %4d | %12s | %s | %s", rows, duration, sql, source))
	}
}

func (l gormLogger) Info(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.Infof(fmt.Sprintf("[Gorm] [info] %s | %s", l.trim(fmt.Sprintf(msg, data...)), utils.FileWithLineNum()))
	}
}

func (l gormLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.Warnf(fmt.Sprintf("[Gorm] [warn] %s | %s", l.trim(fmt.Sprintf(msg, data...)), utils.FileWithLineNum()))
	}
}

func (l gormLogger) Error(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.Errorf(fmt.Sprintf("[Gorm] [error] %s | %s", l.trim(fmt.Sprintf(msg, data...)), utils.FileWithLineNum()))
	}
}
