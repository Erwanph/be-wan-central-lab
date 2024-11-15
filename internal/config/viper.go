package config

import (
	"github.com/spf13/viper"
)

func NewViper() (*viper.Viper, error) {
	config := viper.New()

	config.SetConfigName(".env")
	config.SetConfigType("env")
	config.AddConfigPath(".")
	err := config.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config.AutomaticEnv()
	return config, nil
}
