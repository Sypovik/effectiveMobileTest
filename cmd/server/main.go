package main

import (
	"github.com/Sypovik/effectiveMobileTest/internal/db"
	"github.com/Sypovik/effectiveMobileTest/internal/handlers"
	"github.com/Sypovik/effectiveMobileTest/internal/repository"
	"github.com/Sypovik/effectiveMobileTest/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	d := db.Init()
	r := repository.NewPgPersonRepository(d)
	s := services.NewPersonService(r)
	router := gin.Default()
	handlers.RegisterPersonRoutes(router, *s)
	router.Run()
}
