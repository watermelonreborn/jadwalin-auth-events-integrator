package config

import "github.com/spf13/viper"

type MongoConfig struct {
	ConnectionString string
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
	Redis          RedisConfig
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
	config.Database.Username = viper.GetString("database.username")
	config.Database.Password = viper.GetString("database.password")
	config.Redis.Host = viper.GetString("redis.host")
	config.Redis.Password = viper.GetString("redis.password")
	config.Redis.Database = viper.GetInt("redis.database")

	//TODO: Add another config setups

	return &config, nil
}
