package config

import (
	"github.com/spf13/viper"
	"log"
)

const (
	CfgFileName = "config.yml"
	CfgFileType = "yml"
)

func init() {
	viper.SetConfigName(CfgFileName)
	viper.SetConfigType(CfgFileType)
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
}
