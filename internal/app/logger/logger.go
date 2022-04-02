package logger

import (
	"context"
	"freeSSO/internal/app/config"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

//PgxLogger pgx logger adapter for logrus
type PgxLogger struct {
	l *logrus.Entry
}

func (l *PgxLogger) Log(_ context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
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

func applyLoggerOptions(logger *logrus.Logger, loglevel logrus.Level) {
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	logger.SetLevel(loglevel)
}

//GetNamedLogger returns named logger for diferent modules of app
func GetNamedLoggerWithLevel(name string, loglevel logrus.Level) *logrus.Entry {
	logger := logrus.New()
	applyLoggerOptions(logger, loglevel)
	return logger.WithField("logger", name)
}

func GetNamedLogger(name string) *logrus.Entry {
	conf := config.GetAppConfig()
	var level logrus.Level
	if conf.Debug {
		level = logrus.DebugLevel
	} else {
		level = logrus.InfoLevel
	}
	return GetNamedLoggerWithLevel(name, level)
}
func GetLogger() *logrus.Entry {
	return GetNamedLogger("app")
}

func GetPgxLogger() *PgxLogger {
	logger := logrus.New()
	applyLoggerOptions(logger, logrus.DebugLevel)
	return &PgxLogger{l: logger.WithField("logger", "database")}
}
