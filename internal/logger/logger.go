// logger/logger.go
package logger

import (
	"os"
	"strings"

	"github.com/Sypovik/effectiveMobileTest/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	config := config.LoadConfig()

	// Читаем переменные окружения
	logLevel := strings.ToLower(config.LogLevel)
	usePretty := config.LogPretty

	// Устанавливаем глобальный уровень (по умолчанию — Info)
	level := zerolog.InfoLevel
	switch logLevel {
	case "debug":
		level = zerolog.DebugLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	case "fatal":
		level = zerolog.FatalLevel
	}
	zerolog.SetGlobalLevel(level)

	// Конфигурируем вывод: либо красивый консольный, либо JSON
	var logger zerolog.Logger

	if usePretty {
		writer := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "04:05"}
		logger = zerolog.New(writer).With().Timestamp().Logger()
	} else {
		logger = zerolog.
			New(os.Stdout).
			With().
			Timestamp().
			Logger()
	}

	// Назначаем наш logger как глобальный
	log.Logger = logger
}
