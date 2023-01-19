package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var DBCpm string

func InitConfig() {
	// set default variable
	viper.SetDefault("db.cpm", "sqlserver://sa:yourStrongPassword@localhost:1433?database=dbname")

	// set config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	// viper.AllowEmptyEnv(true)
	// แปลง _ underscore ใน env เป็น . dot notation ใน viper
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// viper.SetEnvKeyReplacer(strings.NewReplacer("_", "."))

	err := viper.ReadInConfig()
	if err != nil {
		zap.L().Warn(fmt.Sprintf("Fatal error config file: %s \n", err))
	}

	DBCpm = viper.GetString("db.cpm")
	// DBCpm = fmt.Sprintf("%#v", viper.AllKeys())

}
