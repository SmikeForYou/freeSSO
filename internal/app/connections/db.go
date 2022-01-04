package connections

import (
	"freeSSO/internal/app/config"
	"freeSSO/internal/app/logger"
	"github.com/jackc/pgx"
	"sync"
)

var pool *pgx.ConnPool = nil

//GetDbConnPool returns db connections pool
func GetDbConnPool() *pgx.ConnPool {
	var once sync.Once
	once.Do(func() {
		appConf := config.GetAppConfig()
		connConf := pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host:                 appConf.DbConfig.Host,
				Port:                 appConf.DbConfig.Port,
				Database:             appConf.DbConfig.Name,
				User:                 appConf.DbConfig.User,
				Password:             appConf.DbConfig.Password,
				TLSConfig:            nil,
				UseFallbackTLS:       false,
				FallbackTLSConfig:    nil,
				Logger:               logger.DbLogger,
				LogLevel:             pgx.LogLevel(appConf.DbConfig.LogLevel),
				Dial:                 nil,
				RuntimeParams:        nil,
				OnNotice:             nil,
				CustomConnInfo:       nil,
				CustomCancel:         nil,
				PreferSimpleProtocol: false,
				TargetSessionAttrs:   "",
			},
			MaxConnections: int(appConf.DbConfig.PoolSize),
			AfterConnect:   nil,
			AcquireTimeout: 0,
		}
		var err error
		pool, err = pgx.NewConnPool(connConf)
		if err != nil {
			logger.AppLogger.Fatal(err)
		}
		logger.AppLogger.Infof("Succesfully established connection to database: %s\n", appConf.DbConfig.ConnStr())
	})
	return pool
}

//CloseDbConnPool closes db connections pool
func CloseDbConnPool() {
	appConf := config.GetAppConfig()
	logger.AppLogger.Infof("Closing db connection pool to database: %s\n", appConf.DbConfig.ConnStr())
	pool.Close()
}
