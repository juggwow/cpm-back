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
	// DBCpm = fmt.Sprintf("%#v", viper.AllKeys())

}
