package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var AppURL string
var AppPort string
var DBCpm string

func InitConfig() {
	viper.SetDefault("app.url", "http://localhost:8000")
	viper.SetDefault("app.port", "8000")
	// set default variable
	viper.SetDefault("db.cpm", "sqlserver://sa:yourStrongPassword@localhost:1433?database=dbname")

	// set config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("Fatal error config file : %s \n", err)
	}

	viper.AutomaticEnv()

	AppURL = viper.GetString("app.url")
	AppPort = viper.GetString("app.port")
	DBCpm = viper.GetString("db.cpm")

}
