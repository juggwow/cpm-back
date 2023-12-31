package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var WebURL string
var WebLoadFileURL string
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

const FILE_SIZE_LIMIT int64 = 52428800

func InitConfig() {
	viper.SetDefault("app.url", "http://localhost:8000")
	viper.SetDefault("app.port", "8000")
	viper.SetDefault("db.rad", "sqlserver://devpool_rad:X1CreIrddfAa5BR4P13resqbUzVGVqop@10.4.34.117:50868?database=RAD")

	viper.SetDefault("web.url", "http://localhost:4200")
	viper.SetDefault("web.loadFile.url", "http://localhost:4200/file")

	viper.SetDefault("auth.callback.url", "http://localhost:4200/contract/31")
	viper.SetDefault("auth.client.id", "CMDC")
	viper.SetDefault("auth.client.secret", "c31bfd34-5de8-4630-a667-9864c02ae455")
	viper.SetDefault("auth.url", "https://sso.pea.co.th/auth/realms/idm")

	viper.SetDefault("auth.jwt.key", "rad-token")
	viper.SetDefault("auth.jwt.secret", "super-secret")
	viper.SetDefault("auth.state", "state")

	viper.SetDefault("storage.endpoint", "minio-api-kolpos.pea.co.th")
	viper.SetDefault("storage.accessKey", "RhoF4o6NbIHzyiei")
	viper.SetDefault("storage.secretKey", "F2epUU6tAeAFBeOB7OGl1DIVaLacmzBc")
	viper.SetDefault("storage.bucketName", "devpool-rad")

	viper.SetDefault("dito.endpoint", "http://172.30.211.224:42/api/pdf-producer/")

	viper.SetDefault("auth.jwt.enabled", true)
	viper.SetDefault("storage.ssl", true)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		zap.L().Warn(fmt.Sprintf("Fatal error config file: %s \n", err))
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	AppURL = viper.GetString("app.url")
	AppPort = viper.GetString("app.port")

	WebURL = viper.GetString("web.url")
	WebLoadFileURL = viper.GetString("web.loadFile.url")

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
