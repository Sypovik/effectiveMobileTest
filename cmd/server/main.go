package main

import (
	_ "github.com/Sypovik/effectiveMobileTest/docs" // docs генерируются Swag CLI
	"github.com/Sypovik/effectiveMobileTest/internal/config"
	"github.com/Sypovik/effectiveMobileTest/internal/db"
	"github.com/Sypovik/effectiveMobileTest/internal/handlers"
	"github.com/Sypovik/effectiveMobileTest/internal/logger"
	"github.com/Sypovik/effectiveMobileTest/internal/middleware"
	"github.com/Sypovik/effectiveMobileTest/internal/repository"
	"github.com/Sypovik/effectiveMobileTest/internal/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Effective Mobile Test API
// @version 1.0
// @description API for managing persons
// @host localhost:8080
// @BasePath /
func main() {
	config := config.LoadConfig()
	port := config.Port

	// инициализируем логгер
	logger.InitLogger()

	// инициализируем БД
	dbConn := db.Init()

	// инициализируем репозиторий
	personRepo := repository.NewPgPersonRepository(dbConn)

	// инициализируем сервис
	personService := services.NewPersonService(personRepo)

	// создаем роутер
	router := gin.Default()

	// добавляем swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// добавляем middleware для логирования
	router.Use(middleware.ZerologContextMiddleware())
	router.Use(middleware.LoggerMiddleware())

	// регистрируем роуты
	handlers.RegisterPersonRoutes(router, *personService)

	// запускаем роутер
	router.Run(":" + port)
}
