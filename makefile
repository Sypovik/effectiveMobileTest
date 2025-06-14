# Makefile для управления Docker Compose и Go приложением (с компиляцией)

# Цели, которые не являются файлами.
.PHONY: build start up down clean restart ps logs help

# Команда по умолчанию
all: help

## up: Запустить все сервисы (app, db) в фоновом режиме
up:
	@echo "🚀 Запуск сервисов (app, db)..."
	docker-compose up -d

## build: Скомпилировать Go-приложение внутри контейнера
build: up
	@echo "🔨 Компиляция Go-приложения внутри контейнера в бинарный файл 'server'..."
	docker-compose exec app go build -v -o ./server ./cmd/server/main.go

## start: Запустить скомпилированное Go-приложение
start: build
	@echo "▶️  Запуск бинарного файла './server' внутри контейнера..."
	docker-compose exec app ./server

## down: Остановить и удалить контейнеры
down:
	@echo "🛑 Остановка контейнеров..."
	docker-compose down

## clean: Остановить контейнеры и удалить все данные (включая тома)
clean:
	@echo "🧹 Полная очистка системы (удаление контейнеров, томов и образов)..."
	docker-compose exec app rm -rf server
	docker-compose down -v --remove-orphans --rmi all

## restart: Перезапустить приложение (полная пересборка и запуск)
restart: down start

## ps: Показать статус запущенных контейнеров
ps:
	docker-compose ps

## logs: Показать логи сервиса 'app' в реальном времени
logs:
	@echo "Просмотр логов..."
	docker-compose logs -f app

## help: Показать это справочное сообщение
help:
	@echo "Доступные команды для make:"
	@echo "------------------------------------------------------------"
	@echo "  up       ->  Запустить сервисы (контейнеры) в фоновом режиме"
	@echo "  build    ->  Скомпилировать Go-приложение внутри контейнера"
	@echo "  start    ->  Запустить скомпилированное приложение (выполняет 'build' перед запуском)"
	@echo "  down     ->  Остановить и удалить контейнеры"
	@echo "  clean    ->  Полностью очистить всё (контейнеры, сети, тома, образы)"
	@echo "  restart  ->  Перезапустить приложение (down + start)"
	@echo "  ps       ->  Показать статус контейнеров"
	@echo "  logs     ->  Показать логи сервиса 'app'"
	@echo "------------------------------------------------------------"