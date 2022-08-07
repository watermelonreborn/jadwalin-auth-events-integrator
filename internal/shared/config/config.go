package config

import "github.com/spf13/viper"

type MongoConfig struct {
	ConnectionString string
	Name             string
	Username         string
	Password         string
}

type RedisConfig struct {
	Host     string
	Password string
	Database int
}

type Config struct {
	HTTPServerPort int
	Database       MongoConfig
	Cache          RedisConfig
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
	config.Database.ConnectionString = viper.GetString("database.connection_string")
	config.Database.Name = viper.GetString("database.name")
	config.Database.Username = viper.GetString("database.username")
	config.Database.Password = viper.GetString("database.password")
	config.Cache.Host = viper.GetString("cache.host")
	config.Cache.Password = viper.GetString("cache.password")
	config.Cache.Database = viper.GetInt("cache.database")

	//TODO: Add another config setups

	return &config, nil
}
