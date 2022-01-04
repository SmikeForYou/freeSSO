package logger

import (
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	"sync"
)

//PgxLogger pgx logger adapter for logrus
type PgxLogger struct {
	l logrus.FieldLogger
}

func (l *PgxLogger) Log(level pgx.LogLevel, msg string, data map[string]interface{}) {
	var logger logrus.FieldLogger
	if data != nil {
		logger = l.l.WithFields(data)
	} else {
		logger = l.l
	}

	switch level {
	case pgx.LogLevelTrace:
		logger.WithField("PGX_LOG_LEVEL", level).Debug(msg)
	case pgx.LogLevelDebug:
		logger.Debug(msg)
	case pgx.LogLevelInfo:
		logger.Info(msg)
	case pgx.LogLevelWarn:
		logger.Warn(msg)
	case pgx.LogLevelError:
		logger.Error(msg)
	default:
		logger.WithField("INVALID_PGX_LOG_LEVEL", level).Error(msg)
	}
}

var AppLogger *logrus.Logger = nil
var DbLogger *PgxLogger = nil

func initAppLogger() *logrus.Logger {
	var once sync.Once
	once.Do(func() {
		AppLogger = logrus.New()
	})
	return AppLogger
}
func initDbLogger() *PgxLogger {
	var once sync.Once
	once.Do(func() {
		DbLogger = &PgxLogger{}
	})
	return DbLogger
}

func init() {
	initAppLogger()
	initDbLogger()
}
