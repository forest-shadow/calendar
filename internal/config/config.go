package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Env      string   `mapstructure:"env"`
	DB       DB       `mapstructure:"db"`
	HTTP     HTTP     `mapstructure:"http"`
	Notifier Notifier `mapstructure:"notifier"`
	Cleaner  Cleaner  `mapstructure:"cleaner"`
	Kafka    Kafka    `mapstructure:"kafka"`
}

type HTTP struct {
	Port int `mapstructure:"port"`
}

type DB struct {
	URI string `mapstructure:"uri"`
}

type Notifier struct {
	LogPath  string `mapstructure:"log_path"`
	Interval string `mapstructure:"interval"`
}

type Cleaner struct {
	Interval string `mapstructure:"interval"`
}

type Kafka struct {
	BootstrapServers string `mapstructure:"bootstrap_servers"`
	Topic            string `mapstructure:"topic"`
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
