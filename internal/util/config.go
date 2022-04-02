package util

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config Estructura para guardar la configuracion global
type Config struct {
	ServerConfig   ServerConfig
	DatabaseConfig DatabaseConfig
	MaxStreak      int `mapstructure:"maxStreak"`
	MaxCount       int `mapstructure:"maxCount"`
}

//ServerConfig Estructura para guardar la configuracion del server
type ServerConfig struct {
	BindAddress string `mapstructure:"bindaddress"`
}

//DatabaseConfig Estructura para guardar la configuracion de la db
type DatabaseConfig struct {
	ClusterModeEnabled bool     `mapstructure:"clustermodeenabled"`
	Addresses          []string `mapstructure:"addresses"`
	PoolSize           int      `mapstructure:"pool_size"`
}

func LoadConfig(path string, file string) (Config, error) {
	var config Config
	var err error
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")
	viper.SetConfigName(file)

	err = viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	var serverConfig ServerConfig
	err = viper.UnmarshalKey("server", &serverConfig)
	if err != nil {
		return config, err
	}

	var databaseConfig DatabaseConfig
	err = viper.UnmarshalKey("database", &databaseConfig)
	if err != nil {
		return config, err
	}

	var maxStreak int
	err = viper.UnmarshalKey("maxStreak", &maxStreak)
	if err != nil {
		return config, err
	}

	var maxCount int
	err = viper.UnmarshalKey("maxCount", &maxCount)
	if err != nil {
		return config, err
	}

	config.ServerConfig = serverConfig
	config.DatabaseConfig = databaseConfig
	config.MaxStreak = maxStreak
	config.MaxCount = maxCount

	fmt.Printf("Configuration loaded.\n")

	return config, err
}
