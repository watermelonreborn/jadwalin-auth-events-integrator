package config

import "github.com/spf13/viper"

type Config struct {
	HTTPServerPort int

	//TODO: Add another config fields
}

func initViper() error {
	viper.SetConfigFile("config.json")

	err := viper.ReadInConfig()
	return err
}

func NewConfig() (*Config, error) {
	var config Config

	if err := initViper(); err != nil {
		return nil, err
	}

	config.HTTPServerPort = viper.GetInt("server.port")

	//TODO: Add another config setups

	return &config, nil
}
