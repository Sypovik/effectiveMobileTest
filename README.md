# API для управления персонами

## Описание
REST API сервис для управления персонами с использованием Gin Gonic и GORM с PostgreSQL.

## Структура проекта

```
effectiveMobileTest/
├── cmd/                 # Пункт входа приложения
│   └── server/
│       └── main.go     # Основной файл приложения
├── internal/           # Внутренний код
│   ├── models/        # Модели данных
│   ├── repository/    # Репозитории
│   ├── services/      # Бизнес-логика
│   ├── handlers/      # HTTP хендлеры
│   └── middleware/    # Middleware
├── migrations/        # SQL миграции
├── docs/              # Документация
└── docker/           # Docker конфигурация
```

## Требования

- Go 1.24.3 или выше
- Docker и Docker Compose
- PostgreSQL 17

## Установка и запуск

### С помощью Docker Compose

1. Установите зависимости:
```bash
make install
```

2. Запустите сервисы:
```bash
docker-compose up -d
```

3. Сервис будет доступен по адресу: http://localhost:8080

### Локально

1. Установите зависимости:
```bash
make install
```

2. Создайте базу данных PostgreSQL:
```sql
CREATE DATABASE people;
```

3. Настройте переменные окружения в файле `.env`:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=people
PORT=8080
LOG_LEVEL=debug
LOG_PRETTY=true
```

4. Запустите миграции:
```bash
make migrate
```

5. Запустите приложение:
```bash
make run
```

## API endpoints

### Создание персоны
```http
POST /people
Content-Type: application/json

{
    "name": "Имя",
    "surname": "Фамилия",
    "patronymic": "Отчество",
    "age": 30,
    "gender": "male/female",
    "country": "RU"
}
```

### Получение списка персон
```http
GET /people

# Фильтрация
?name=Имя
?surname=Фамилия
?patronymic=Отчество
?gender=male/female
?country=RU
?age=30
```

### Получение персоны по ID
```http
GET /people/{id}
```

### Обновление персоны
```http
PUT /people/{id}
Content-Type: application/json

{
    "name": "Имя",
    "surname": "Фамилия",
    "patronymic": "Отчество",
    "age": 30,
    "gender": "male/female",
    "country": "RU"
}
```

### Удаление персоны
```http
DELETE /people/{id}
```

## Логирование

Система использует zerolog для логирования с двумя уровнями:

- INFO: URL запроса и статус ответа
- DEBUG: Полная информация о запросе (метод, путь, параметры, статус, время выполнения)

Настройки логирования:
```
LOG_LEVEL=debug
LOG_PRETTY=true
```

## Тестирование

Для запуска тестов:
```bash
make test
```

## Миграции

Для выполнения миграций:
```bash
make migrate
```