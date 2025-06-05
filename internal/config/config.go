package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// DB_HOST
// DB_PORT
// DB_USER
// DB_PASSWORD
// DB_NAME
// PORT

type Config struct {
	Port       string `env:"PORT"`
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`
	DBSslmode  string `env:"DB_SSLMODE"`
}

func LoadConfig() *Config {
	fmt.Println("asdasd")
	configPath := ".env"
	var cfg Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", err)
	}

	err := cleanenv.ReadConfig(configPath, &cfg)
	if os.IsNotExist(err) {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
