package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// LoggerMiddleware логирует HTTP запросы
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var url string
		if c.Request.URL.RawQuery == "" {
			url = c.Request.URL.Scheme + "://" + c.Request.Host + c.Request.URL.Path
		} else {
			url = c.Request.URL.Scheme + "://" + c.Request.Host + c.Request.URL.Path + "?" + c.Request.URL.RawQuery
		}
		start := time.Now()
		logger := log.Ctx(c.Request.Context())

		// Логируем входящий запрос (INFO)
		logger.
			Info().
			Str("url", url).
			Msg("HTTP запрос получен")

		// Продолжаем обработку запроса
		c.Next()

		// Логируем ответ (INFO)
		logger.
			Info().
			Str("url", url).
			Int("status", c.Writer.Status()).
			Msg("HTTP запрос завершен")

		// Логируем детали запроса (DEBUG)
		logger.
			Debug().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("query", c.Request.URL.RawQuery).
			Str("client_ip", c.ClientIP()).
			Int("status", c.Writer.Status()).
			Int64("latency_ms", time.Since(start).Milliseconds()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("Детали HTTP запроса")
	}
}
