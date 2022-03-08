package connections

import (
	"freeSSO/internal/app/config"
	"freeSSO/internal/app/logger"

	"github.com/jackc/pgx"
)

var pool *pgx.ConnPool = nil
var log = logger.GetLogger()

//GetDbConnPool returns db connections pool
func GetDbConnPool() *pgx.ConnPool {
	if pool == nil {
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
				Logger:               logger.GetPgxLogger(),
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
			log.Fatal(err)
		}
		log.Infof("Succesfully established connection to database: %s\n", appConf.DbConfig.ConnStr())
	}
	return pool
}

//CloseDbConnPool closes db connections pool
func CloseDbConnPool() {
	appConf := config.GetAppConfig()
	log.Infof("Closing db connection pool to database: %s\n", appConf.DbConfig.ConnStr())
	pool.Close()
}
