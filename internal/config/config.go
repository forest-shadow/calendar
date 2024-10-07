package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DB       DBConfig
	HTTP     HTTPConfig
}
	
type HTTPConfig struct {
	Port int `env:"HTTP_PORT" envDefault:"8089"`
}

type DBConfig struct {
	URI string `env:"DB_URI" envDefault:"postgresql://postgres:password@localhost:5555/auth"`
}

func envLoader() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}
}

func GetConfig() (*Config, error) {
	envLoader()

	cfg := &Config{}
    if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("load config from env: %w", err)
	}

	return cfg, nil
}