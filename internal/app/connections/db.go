package connections

import (
	"context"
	"freeSSO/internal/app/config"
	"freeSSO/internal/app/logger"

	"github.com/jackc/pgx/v4"
)

var pool *pgx.Conn = nil
var log = logger.GetLogger()

//GetDbConnPool returns db connections pool
func GetDbConnPool() *pgx.Conn {
	if pool == nil {
		appConf := config.GetAppConfig()
		connConf, _ := pgx.ParseConfig(appConf.DbConfig.ConnStr())
		connConf.Logger = logger.GetPgxLogger()
		connConf.LogLevel = pgx.LogLevel(appConf.DbConfig.LogLevel)
		var err error
		pool, err = pgx.ConnectConfig(context.Background(), connConf)
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
	pool.Close(context.TODO())
}
