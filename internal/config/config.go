package config

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	env := os.Getenv("APP_ENV")
	if env == "" {
		log.Printf("APP_ENV is not set, defaulting to 'local'")
		env = "local"
	}

	viper.SetConfigName("env." + env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("error reading config file: %v", err)
	}

	viper.SetEnvPrefix("app")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &config, nil
}
