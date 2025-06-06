# **Go REST API для управления персонами с обогащением данных**

Это RESTful API, разработанное на Go, для управления информацией о персонах. При создании и обновлении персон, данные обогащаются информацией о возрасте, поле и национальности через внешние API. Проект включает в себя пагинацию, фильтрацию, логирование, а также документацию Swagger.

## **Структура проекта**
```sh
.  
├── cmd  
│   └── server  
│       └── main.go           # Точка входа в приложение  
├── docker-compose.yaml     # Конфигурация Docker Compose для развертывания  
├── Dockerfile              # Dockerfile для сборки образа приложения  
├── docs  
│   ├── docs.go             # Сгенерированные Go Swagger документы  
│   ├── swagger.json        # Сгенерированная спецификация Swagger в JSON  
│   └── swagger.yaml        # Сгенерированная спецификация Swagger в YAML  
├── go.mod                  # Модуль Go  
├── go.sum                  # Сумма Go модулей  
├── internal  
│   ├── config  
│   │   └── config.go       # Загрузка конфигурации приложения из .env  
│   ├── db  
│   │   └── db.go           # Инициализация подключения к базе данных  
│   ├── dto  
│   │   ├── person_request.go   # DTO для запросов (создание, обновление)  
│   │   └── person_response.go  # DTO для ответов (получение, список)  
│   ├── handlers  
│   │   └── person_handler.go   # Обработчики HTTP-запросов  
│   ├── logger  
│   │   └── logger.go       # Инициализация и настройка zerolog  
│   ├── middleware  
│   │   ├── logger_middleware.go         # Middleware для логирования HTTP запросов  
│   │   └── zerologContextMiddleware.go  # Middleware для добавления zerolog в контекст  
│   ├── models  
│   │   └── person.go           # Модель данных Person  
│   ├── repository  
│   │   ├── person_repository.go    # Интерфейс репозитория для Person  
│   │   └── pg_person.go            # Реализация репозитория Person для PostgreSQL  
│   └── services  
│       ├── person_enricher.go      # Логика обогащения данных персоны  
│       └── person_service.go       # Бизнес-логика для Person  
├── makefile                # Makefile для автоматизации задач  
│  
├── README.md               # Данный файл  
└── TZ.md                   # Дополнительная информация по техническому заданию
```

## **Функциональность**

API предоставляет следующие возможности по управлению персонами:

1. **Создание персоны (POST /people)**:  
   - Принимает name, surname (обязательные) и patronymic (опционально).  
   - Автоматически обогащает данные age, gender, country для переданного name с использованием внешних API: api.agify.io, api.genderize.io, api.nationalize.io.  
   - Сохраняет персону в базе данных.  
2. **Получение персоны по ID (GET /people/{id})**:  
   - Возвращает полную информацию о персоне по её уникальному идентификатору.  
3. **Обновление персоны (PUT /people/{id})**:  
   - Позволяет обновлять поля name, surname, patronymic, age, gender, country.  
   - Если поле name изменяется, происходит повторное обогащение данных age, gender, country.  
4. **Удаление персоны (DELETE /people/{id})**:  
   - Удаляет запись о персоне из базы данных по её ID.  
5. **Получение списка персон (GET /people)**:  
   - Поддерживает пагинацию (параметры page и size).  
   - Поддерживает фильтрацию по полям: name, surname, patronymic, gender, country, age.  
   - Возвращает список персон и общее количество записей, соответствующих фильтрам.

## **Примеры запросов с cURL**

Ниже представлены примеры HTTP-запросов к API с использованием утилиты cURL. Предполагается, что приложение запущено и доступно по адресу http://localhost:8080.

### **1. Создание персоны (POST /people)**

Создание новой персоны с обязательными полями name и surname. Поля age, gender, country будут обогащены автоматически.
```
curl -X POST http://localhost:8080/people   
-H "Content-Type: application/json"   
-d '{  
    "name": "Иван",  
    "surname": "Иванов",  
    "patronymic": "Иванович"  
}'
```
Пример ответа:
```
{  
    "id": 1,  
    "name": "Иван",  
    "surname": "Иванов",  
    "patronymic": "Иванович",  
    "age": 42,  
    "gender": "male",  
    "country": "RU",  
    "created_at": "2025-06-06T13:54:59Z",  
    "updated_at": "2025-06-06T13:54:59Z"  
}
```
### **2. Получение персоны по ID (GET /people/{id})**

Получение информации о персоне с ID = 1.
```
curl -X GET http://localhost:8080/people/1
```
Пример ответа:
```
{  
    "id": 1,  
    "name": "Иван",  
    "surname": "Иванов",  
    "patronymic": "Иванович",  
    "age": 42,  
    "gender": "male",  
    "country": "RU",  
    "created_at": "2025-06-06T13:54:59Z",  
    "updated_at": "2025-06-06T13:54:59Z"  
}
```
### **3. Обновление персоны (PUT /people/{id})**

Обновление имени персоны с ID = 1. Если поле name изменяется, данные age, gender, country будут повторно обогащены.
```
curl -X PUT http://localhost:8080/people/1   
-H "Content-Type: application/json"   
-d '{  
    "name": "Алексей"  
}'
```
Пример ответа:
```
{  
    "id": 1,  
    "name": "Алексей",  
    "surname": "Иванов",  
    "patronymic": "Иванович",  
    "age": 35,  
    "gender": "male",  
    "country": "RU",  
    "created_at": "2025-06-06T13:54:59Z",  
    "updated_at": "2025-06-06T13:55:00Z"  
}
```
### **4. Получение списка персон (GET /people)**

Получение списка всех персон с пагинацией (10 записей на первой странице).
```
curl -X GET "http://localhost:8080/people?page=1&size=10"
```
Пример ответа:
```json
{  
    "total": 8,  
    "page": 1,  
    "size": 10,  
    "data": [  
        {  
            "id": 1,  
            "name": "Nikita",  
            "surname": "Ushakov",  
            "patronymic": "Vasilevich",  
            "age": 45,  
            "gender": "male",  
            "country": "RU",  
            "created_at": "2025-06-06T13:54:59Z",  
            "updated_at": "2025-06-06T13:54:59Z"  
        },  
        {  
            "id": 2,  
            "name": "Anton",  
            "surname": "Ushakov",  
            "patronymic": "Vasilevich",  
            "age": 58,  
            "gender": "male",  
            "country": "RO",  
            "created_at": "2025-06-06T13:55:00Z",  
            "updated_at": "2025-06-06T13:55:00Z"  
        }  
    ]  
}
```
Фильтрация по имени и полу:

``` 
curl -X GET "http://localhost:8080/people?name=Иван&gender=male"
```
### **5. Удаление персоны (DELETE /people/{id})**

Удаление персоны с **ID = 1**

``` 
curl -X DELETE http://localhost:8080/people/1
```
При успешном удалении будет возвращен статус 204 No Content.

## **Используемые технологии**

- **Go (Golang)**: Основной язык разработки.  
- **Gin Gonic**: Легковесный HTTP веб-фреймворк для Go.  
- **GORM**: ORM-библиотека для Go, используемая для взаимодействия с PostgreSQL.  
- **PostgreSQL**: Реляционная база данных для хранения информации о персонах.  
- **Zerolog**: Быстрая и простая в использовании библиотека логирования.  
- **Godotenv**: Для загрузки переменных окружения из файла .env.  
- **Cleanenv**: Для чтения конфигурации приложения из переменных окружения.  
- **Go Swagger**: Для автоматической генерации документации API.  
- **Docker & Docker Compose**: Для контейнеризации и оркестрации приложения и базы данных.

## **Запуск проекта**

### **1. Предварительные требования**

Убедитесь, что у вас установлены:

* [Go](https://golang.org/doc/install) (версия 1.20 или выше)  
* [Docker](https://docs.docker.com/get-docker/)  
* [Docker Compose](https://docs.docker.com/compose/install/)

### **2. Клонирование репозитория**
```sh
git clone git@github.com:Sypovik/effectiveMobileTest.git  
cd effectiveMobileTest
```

### **3. Настройка переменных окружения**

Файл .env в корневой директории проекта со следующим содержимым:
```sh
DB_HOST=db  
DB_PORT=5432  
DB_USER=postgres  
DB_PASSWORD=postgres  
DB_NAME=people  
DB_SSLMODE=disable  
PORT=8080   
LOG_LEVEL=debug #info
LOG_PRETTY=true
```

### **4. Запуск с помощью Docker Compose**

Это предпочтительный способ запуска, так как он автоматически поднимет базу данных PostgreSQL, применит миграции и запустит API.
```sh
docker-compose up --build
```

После запуска приложение будет доступно по адресу http://localhost:8080.

### **5. Запуск без Docker (для разработки)**

Если вы хотите запустить приложение без Docker (только Go часть), убедитесь, что у вас есть запущенный PostgreSQL сервер и создана база данных person_db с пользователем user и паролем password.

1. **Установка зависимостей:**  
```sh
go mod tidy
```

2. **Запуск приложения:**  
```sh
go run cmd/server/main.go
```

## **Документация API (Swagger)**

Документация API доступна по адресу http://localhost:8080/swagger/index.html после запуска приложения.

Для обновления документации Swagger вручную (после изменений в коде с // @ аннотациями):
```sh
# Установите swag, если ещё не установлен:  
go install github.com/swaggo/swag/cmd/swag@latest  
swag init -g cmd/server/main.go
```

## **Логирование**

Приложение использует zerolog для логирования. Уровень логирования и формат вывода (pretty-print для консоли или JSON) можно настроить через переменные окружения LOG_LEVEL и LOG_PRETTY в файле .env.

Примеры уровней логирования: debug, info, warn, error, fatal.

## **Makefile**

Проект содержит Makefile для автоматизации некоторых задач:

- `make up`: Запустить сервисы (контейнеры app, db) в фоновом режиме.  
- `make build`: Скомпилировать Go-приложение внутри контейнера в бинарный файл server.  
- `make start`: Запустить скомпилированное приложение (выполняет build перед запуском).  
- `make down`: Остановить и удалить контейнеры.  
- `make clean`: Полностью очистить всё (контейнеры, сети, тома, образы, а также удаляет бинарник server из контейнера app).  
- `make restart`: Перезапустить приложение (останавливает, затем запускает).  
- `make ps`: Показать статус запущенных контейнеров.  
- `make logs`: Показать логи сервиса app в реальном времени.  
- `make help`: Показать справочное сообщение со всеми доступными командами.

## **Дальнейшие улучшения / Замечания**

- Добавление аутентификации и авторизации.  
- Более продвинутая обработка ошибок и кастомные ответы для клиентской стороны.  
- Реализация rate limiting для внешних API.  
- Добавление юнит- и интеграционных тестов.  
- Более полное покрытие логами разного уровня.  
- Внедрение кэширования для обогащенных данных.  
- Улучшение удобства использования документации API (например, предоставление Postman Collection).