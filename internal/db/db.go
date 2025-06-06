package db

import (
	"fmt"
	"log"
	"time"

	"github.com/Sypovik/effectiveMobileTest/internal/config"
	"github.com/Sypovik/effectiveMobileTest/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// var DB *gorm.DB

// Init устанавливает соединение с БД и применяет миграции.
// Если вы хотите откладывать применение миграций до отдельного шага (например, через makefile),
// можно вынести вызов applyMigrations в main.go.
func Init() *gorm.DB {
	// 1. Считать параметры из окружения
	config := config.LoadConfig()
	host := config.DBHost
	port := config.DBPort
	password := config.DBPassword
	dbname := config.DBName
	user := config.DBUser
	sslmode := config.DBSslmode // обычно "disable", можно вынести в .env

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("DB: отсутствуют обязательные переменные окружения (DB_HOST/DB_PORT/DB_USER/DB_PASSWORD/DB_NAME)")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		host, user, password, dbname, port, sslmode,
	)

	// 2. Открыть соединение через GORM
	// DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
	// 	Logger: logger.Default.LogMode(logger.Info),
	// })
	DB, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatalf("DB: не удалось подключиться к базе данных: %v", err)
	}
	DB.AutoMigrate(&models.Person{})
	// 3. Настроить параметры «низкоуровневого» sql.DB
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("DB: не удалось получить объект sql.DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	log.Println("DB: подключение установлено")

	return DB
}
