package conn

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

// GormLogger struct
type GormLogger struct {
	Logger *logrus.Logger
}

func NewGormLogger(logger *logrus.Logger) *GormLogger {
	return &GormLogger{Logger: logger}
}

// Print - Log Formatter
func (g *GormLogger) Print(v ...interface{}) {
	// v[0]: type; v[1]: src; v[2]: duration; v[3]: sql; v[4]: values; v[5]: rows
	switch v[0] {
	case "sql":
		g.Logger.WithFields(logrus.Fields{
			"Module":   "gorm",
			"Type":     "sql",
			"Source":   v[1],
			"Duration": v[2],
			"SQL":      v[3],
			"Values":   v[4],
			"Rows":     v[5],
		}).Info(fmt.Sprintf("[Gorm] rows: %3d | %10s | %s", v[5], v[2], v[3]))
	case "log":
		g.Logger.WithFields(logrus.Fields{
			"Module": "gorm",
			"Type":   "log",
		}).Print(fmt.Sprintf("[Gorm] %s", v[2]))
	}
}
