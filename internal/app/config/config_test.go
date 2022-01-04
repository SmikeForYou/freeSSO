package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetAppConfig(t *testing.T) {
	testKeys := []string{
		"MEM_DB_HOST",
		"MEM_DB_PORT",
		"MEM_DB_PASSWORD",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
		"DB_USER",
		"DB_PASSWORD",
		"DB_POOL_SIZE",
		"DB_SSL_MODE",
		"DB_SSL_ROOT_CERT_PATH",
		"DB_SSL_KEY_PATH",
		"DB_SSL_CERT_PATH",
		"GRPC_SERVER_PORT",
		"HTTP_SERVER_PORT",
	}
	oldValues := make(map[string]string, len(testKeys))
	for _, key := range testKeys {
		oldValues[key] = os.Getenv(key)
	}
	_ = os.Setenv("MEM_DB_HOST", "memdbhost")
	_ = os.Setenv("MEM_DB_PORT", "6432")
	_ = os.Setenv("MEM_DB_PASSWORD", "password")
	_ = os.Setenv("DB_HOST", "dbhost")
	_ = os.Setenv("DB_PORT", "27017")
	_ = os.Setenv("DB_NAME", "db")
	_ = os.Setenv("DB_USER", "dbuser")
	_ = os.Setenv("DB_PASSWORD", "dbpass")
	_ = os.Setenv("DB_POOL_SIZE", "15")
	_ = os.Setenv("DB_SSL_MODE", "true")
	_ = os.Setenv("DB_SSL_ROOT_CERT_PATH", "/path")
	_ = os.Setenv("DB_SSL_KEY_PATH", "/path")
	_ = os.Setenv("DB_SSL_CERT_PATH", "/path")
	_ = os.Setenv("HTTP_SERVER_PORT", "7772")
	_ = os.Setenv("GRPC_SERVER_PORT", "7771")

	appConfig = GetAppConfig()
	assert.Equal(t, "memdbhost", appConfig.MemDbConfig.Host)
	assert.Equal(t, uint16(6432), appConfig.MemDbConfig.Port)
	assert.Equal(t, "password", appConfig.MemDbConfig.Password)
	assert.Equal(t, "dbhost", appConfig.DbConfig.Host)
	assert.Equal(t, uint16(27017), appConfig.DbConfig.Port)
	assert.Equal(t, "db", appConfig.DbConfig.Name)
	assert.Equal(t, "dbuser", appConfig.DbConfig.User)
	assert.Equal(t, "dbpass", appConfig.DbConfig.Password)
	assert.Equal(t, uint16(15), appConfig.DbConfig.PoolSize)
	assert.Equal(t, "true", appConfig.DbConfig.SslMode)
	assert.Equal(t, "/path", appConfig.DbConfig.SslRootCertPath)
	assert.Equal(t, "/path", appConfig.DbConfig.SslKeyPath)
	assert.Equal(t, "/path", appConfig.DbConfig.SslCertPath)
	assert.Equal(t, uint16(27017), appConfig.DbConfig.Port)
	assert.Equal(t, uint16(7771), appConfig.GrpcServerConfig.Port)
	assert.Equal(t, uint16(7772), appConfig.HttpServerConfig.Port)
	assert.Equal(t, "host=dbhost port=27017 user=dbuser password=dbpass dbname=db sslmode=disable", appConfig.DbConfig.ConnStr())
	anotherAppConfig := GetAppConfig()
	assert.Equal(t, appConfig, anotherAppConfig)
	for key, value := range oldValues {
		_ = os.Setenv(key, value)
	}

}
