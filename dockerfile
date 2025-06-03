# Используем легковесный образ Go
FROM golang:1.24-alpine

# Установка необходимых утилит

RUN go get github.com/go-pg/pg/v10
RUN go get -u github.com/gin-gonic/gin

# Рабочая директория внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Копируем всё остальное
COPY . .

# Собираем приложение
# RUN go build -o app ./cmd/server

# Указываем порт
EXPOSE 8080

# Команда запуска
# CMD ["./app"]
