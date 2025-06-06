package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// LoggerMiddleware добавляет zerolog.Logger в контекст каждого запроса
func ZerologContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Создаём новый контекст с логером
		ctx := log.Logger.WithContext(c.Request.Context())
		// Вписываем в запрос новый контекст с логером
		c.Request = c.Request.WithContext(ctx)

		// Далее вызываем следующий обработчик
		c.Next()
	}
}
