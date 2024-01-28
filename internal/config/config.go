package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type DBConfig struct {
	User     string   `env:"DB_USER,required"`
	Password string   `env:"DB_PASSWORD,required"`
	Host     string   `env:"DB_HOST,required"`
	Port     int      `env:"DB_PORT,required"`
	DBName   string   `env:"DB_NAME,required"`
	Tables   []string `env:"DB_TABLES,required"`
}

type ServerConfig struct {
	Address     string        `env:"SERVER_ADDRESS,required"`
	Timeout     time.Duration `env:"SERVER_TIMEOUT"`
	IdleTimeout time.Duration `env:"SERVER_IDLE_TIMEOUT"`
}

type Config struct {
	Database DBConfig
	Server   ServerConfig
	Env      string `env:"ENV_LEVEL"`
}

func ConfigInit() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("empty config path")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", configPath)
	}

	err := godotenv.Load(configPath)
	if err != nil {
		return nil, err
	}

	cfg := Config{}

	err = env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
