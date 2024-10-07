package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DB       DB
	HTTP     HTTP
}
	
type HTTP struct {
	Port int `env:"HTTP_PORT" envDefault:"8089"`
}

type DB struct {
	URI string `env:"DB_URI" envDefault:"postgresql://postgres:password@localhost:5555/auth"`
}

func envLoader() {
	// Get the environment from GO_ENV, defaulting to 'local' if not set
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "local"
	}
	
	// Determine the path to the appropriate .env file
	envPath := filepath.Join(".", fmt.Sprintf(".env.%s", env))
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("No %v file found", envPath)
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