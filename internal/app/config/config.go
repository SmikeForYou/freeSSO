package config

import (
	"fmt"
	"freeSSO/internal/app/logger"
	"freeSSO/internal/app/utils"
	"sync"
)

type MemDbConfig struct {
	Host     string `env:"MEM_DB_HOST"`
	Port     uint16 `env:"MEM_DB_PORT"`
	Password string `env:"MEM_DB_PASSWORD"`
}

type DbConfig struct {
	Host            string `env:"DB_HOST"`
	Port            uint16 `env:"DB_PORT"`
	Name            string `env:"DB_NAME"`
	User            string `env:"DB_USER"`
	Password        string `env:"DB_PASSWORD"`
	PoolSize        uint16 `env:"DB_POOL_SIZE"`
	LogLevel        uint16 `env:"DB_LOG_LEVEL"`
	SslMode         string `env:"DB_SSL_MODE"`
	SslRootCertPath string `env:"DB_SSL_ROOT_CERT_PATH"`
	SslKeyPath      string `env:"DB_SSL_KEY_PATH"`
	SslCertPath     string `env:"DB_SSL_CERT_PATH"`
}

func (dbConf DbConfig) ConnStr() string {
	//TODO: implement ssl modes for connection
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.Name)
}

type GrpcServerConfig struct {
	Port uint16 `env:"GRPC_SERVER_PORT"`
}

type HttpServerConfig struct {
	Port uint16 `env:"HTTP_SERVER_PORT"`
}

type ActionLoggerConfig struct {
}

type AppConfig struct {
	SecretKey string `env:"APP_SECRET_KEY"`
	Debug     bool   `env:"APP_DEBUG"`
	DbConfig
	MemDbConfig
	HttpServerConfig
	GrpcServerConfig
	ActionLoggerConfig
}

var appConfig AppConfig

// GetAppConfig returns application config
func GetAppConfig() AppConfig {
	var once sync.Once
	once.Do(func() {
		//init base app values
		err := utils.InitStructFromEnv(&appConfig)
		if err != nil {
			logger.AppLogger.Fatalln(err)
		}
		//init db config
		var dbConf DbConfig
		err = utils.InitStructFromEnv(&dbConf)
		if err != nil {
			logger.AppLogger.Fatalln(err)
		}
		appConfig.DbConfig = dbConf
		//init in memory db config
		var memDbConf MemDbConfig
		err = utils.InitStructFromEnv(&memDbConf)
		if err != nil {
			logger.AppLogger.Fatalln(err)
		}
		appConfig.MemDbConfig = memDbConf
		//init grpc server config
		var grpcServerConfig GrpcServerConfig
		err = utils.InitStructFromEnv(&grpcServerConfig)
		if err != nil {
			logger.AppLogger.Fatalln(err)
		}
		appConfig.GrpcServerConfig = grpcServerConfig
		//init http server config
		var httpServerConfig HttpServerConfig
		err = utils.InitStructFromEnv(&httpServerConfig)
		if err != nil {
			logger.AppLogger.Fatalln(err)
		}
		appConfig.HttpServerConfig = httpServerConfig
		//init action logger config
		var actionLoggerConfig ActionLoggerConfig
		err = utils.InitStructFromEnv(&actionLoggerConfig)
		if err != nil {
			logger.AppLogger.Fatalln(err)
		}
		appConfig.ActionLoggerConfig = actionLoggerConfig
		logger.AppLogger.Infof("Running app with conf %+v\n", appConfig)
	})
	return appConfig
}
