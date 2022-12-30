package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var DBCpm string

func InitConfig() {
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

	DBCpm = viper.GetString("db.cpm")

}
