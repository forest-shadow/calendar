package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DB   DB   `mapstructure:"db"`
	HTTP HTTP `mapstructure:"http"`
}

type HTTP struct {
	Port int `mapstructure:"port"`
}

type DB struct {
	URI string `mapstructure:"uri"`
}

func GetConfig() (*Config, error) {
	env := os.Getenv("GO_ENV")
	if env == "" {
		log.Printf("GO_ENV is not set, defaulting to 'local'")
		env = "local"
	}

	viper.SetConfigName("env." + env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &config, nil
}
