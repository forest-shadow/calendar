package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"log"
	path "github.com/forest-shadow/calendar/internal/constants"
)

type Config struct {
	HTTP     HTTPConfig
}
	
type HTTPConfig struct {
	Port int `env:"HTTP_PORT" envDefault:"8089"`
}

func envLoader() {
	if err := godotenv.Load(path.EnvFile); err != nil {
		log.Println("No .env file found")
	}
}


func Parse() (*Config, error) {
	envLoader()

	cfg := &Config{}
    if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("load config from env: %w", err)
	}

	return cfg, nil
}