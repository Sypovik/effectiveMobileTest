package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	// Port - порт, на котором будет запущен HTTP сервер
	// По умолчанию: 8080
	Port string `env:"PORT"`

	// DBHost - хост PostgreSQL базы данных
	// По умолчанию: localhost
	DBHost string `env:"DB_HOST"`

	// DBPort - порт PostgreSQL базы данных
	// По умолчанию: 5432
	DBPort string `env:"DB_PORT"`

	// DBUser - имя пользователя для подключения к базе данных
	DBUser string `env:"DB_USER"`

	// DBPassword - пароль для подключения к базе данных
	DBPassword string `env:"DB_PASSWORD"`

	// DBName - имя базы данных
	DBName string `env:"DB_NAME"`

	// DBSslmode - режим SSL соединения с базой данных
	// Возможные значения: disable, require, prefer, verify-ca, verify-full
	DBSslmode string `env:"DB_SSLMODE"`

	// LogLevel - уровень логирования
	LogLevel string `env:"LOG_LEVEL"`

	// LogPretty - использовать красивый формат вывода логов
	LogPretty bool `env:"LOG_PRETTY"`
}

// LoadConfig загружает конфигурацию из файла .env
// Если файл не найден или есть ошибки при чтении - приложение падает с ошибкой
func LoadConfig() *Config {
	fmt.Println("Загрузка конфигурации...")
	configPath := ".env"
	var cfg Config

	// Проверяем существование файла конфигурации
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Файл конфигурации не найден: %s", err)
	}

	// Парсим конфигурацию из файла
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Ошибка при чтении конфигурации: %s", err)
	}

	return &cfg
}
