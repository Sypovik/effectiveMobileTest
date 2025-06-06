package main

import (
	_ "github.com/Sypovik/effectiveMobileTest/docs" // docs генерируются Swag CLI
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

	logger.InitLogger()

	d := db.Init()
	r := repository.NewPgPersonRepository(d)
	s := services.NewPersonService(r)
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(middleware.ZerologContextMiddleware())
	router.Use(middleware.LoggerMiddleware())
	handlers.RegisterPersonRoutes(router, *s)
	router.Run()
}
