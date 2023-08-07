package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var AppURL string
var AppPort string
var DBCpm string
var AuthCallbackURL string
var AuthClientID string
var AuthClientSecret string
var AuthJWTEnabled bool
var AuthJWTExpiredDuration = time.Duration(4) * time.Hour
var AuthJWTExpiredRefreshDuration = time.Duration(8) * time.Hour
var AuthJWTKey string
var AuthJWTSecret string
var AuthState = xid.New().String()
var AuthURL string
var StorageSSL bool
var StorageEndpoint string
var StorageAccessKey string
var StorageSecretKey string
var StorageBucketName string
var DitoApi string
var DBCon string

func InitConfig() {
	viper.SetDefault("app.url", "http://localhost:8000")
	viper.SetDefault("app.port", "8000")

	viper.SetDefault("web.url", "http://localhost:4200/")
	viper.SetDefault("web.fileAttachment.url", "http://localhost:4200/file/")

	viper.SetDefault("auth.callback.url", "http://localhost:8000/auth/callback")

	viper.SetDefault("auth.jwt.enabled", true)

	viper.SetDefault("storage.ssl", true)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		zap.L().Warn(fmt.Sprintf("Fatal error config file: %s \n", err))
		// DBCon = fmt.Sprintf("Fatal error config file: %s \n", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	DBCon = "env : " + viper.GetString("db.rad")
	// DBCon = DBCon + "\nsecretKey : " + viper.GetString("storage.secretKey")
	// DBCon = DBCon + "\naccessKey : " + viper.GetString("storage.accessKey")

	AppURL = viper.GetString("app.url")
	AppPort = viper.GetString("app.port")

	DBCpm = viper.GetString("db.rad")

	StorageSSL = viper.GetBool("storage.ssl")
	StorageEndpoint = viper.GetString("storage.endpoint")
	StorageAccessKey = viper.GetString("storage.accessKey")
	StorageSecretKey = viper.GetString("storage.secretKey")
	StorageBucketName = viper.GetString("storage.bucketName")

	AuthCallbackURL = viper.GetString("auth.callback.url")
	AuthClientID = viper.GetString("auth.client.id")
	AuthClientSecret = viper.GetString("auth.client.secret")
	AuthJWTEnabled = viper.GetBool("auth.jwt.enabled")
	AuthJWTKey = viper.GetString("auth.jwt.key")
	AuthJWTSecret = viper.GetString("auth.jwt.secret")
	AuthURL = viper.GetString("auth.url")

	DitoApi = viper.GetString("dito.endpoint")

}
