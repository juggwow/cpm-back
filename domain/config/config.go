package config

import (
	"fmt"
	"strings"

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
var AuthJWTKey string
var AuthJWTSecret string
var AuthState string
var AuthURL string
var StorageSSL bool
var StorageEndpoint string
var StorageAccessKey string
var StorageSecretKey string
var StorageBucketName string

func InitConfig() {
	viper.SetDefault("app.url", "http://localhost:8000")
	viper.SetDefault("app.port", "8000")
	// set default variable
	viper.SetDefault("db.rad", "sqlserver://devpool_rad:X1CreIrddfAa5BR4P13resqbUzVGVqop@10.4.34.117:50868?database=RAD")

	viper.SetDefault("auth.callback.url", "http://localhost:3000")
	viper.SetDefault("auth.client.id", "client_id")
	viper.SetDefault("auth.client.secret", "client_secret")
	viper.SetDefault("auth.jwt.enabled", true)
	viper.SetDefault("auth.jwt.key", "cmdc-token")
	viper.SetDefault("auth.jwt.secret", "super-secret")
	viper.SetDefault("auth.state", "state")
	viper.SetDefault("auth.url", "https://sso.pea.co.th/auth/realms/idm")

	viper.SetDefault("storage.ssl", true)
	viper.SetDefault("storage.endpoint", "minio-api-kolpos.pea.co.th")
	// viper.SetDefault("storage.accessKey", "RhoF4o6NbIHzyiei")
	// viper.SetDefault("storage.secretKey", "F2epUU6tAeAFBeOB7OGl1DIVaLacmzBc")
	viper.SetDefault("storage.accessKey", "devpool-rad")
	viper.SetDefault("storage.secretKey", "YHWdKG6zePjiOGNEluK7oE3msPn50HCN")
	viper.SetDefault("storage.bucketName", "devpool-rad")

	// set config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	// viper.AllowEmptyEnv(true)
	// แปลง _ underscore ใน env เป็น . dot notation ใน viper
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		zap.L().Warn(fmt.Sprintf("Fatal error config file: %s \n", err))
	}

	AppURL = viper.GetString("app.url")
	AppPort = viper.GetString("app.port")
	DBCpm = viper.GetString("db.rad")
	StorageSSL = viper.GetBool("storage.ssl")
	StorageEndpoint = viper.GetString("storage.endpoint")
	StorageAccessKey = viper.GetString("storage.accessKey")
	StorageSecretKey = viper.GetString("storage.secretKey")
	StorageBucketName = viper.GetString("storage.bucketName")
	// DBCpm = fmt.Sprintf("%#v", viper.AllKeys())

}
