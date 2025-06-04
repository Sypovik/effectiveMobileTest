# Используем легковесный образ Go
FROM golang:1.24-alpine

# Установка необходимых утилит

WORKDIR /app
COPY go.mod ./


RUN go get -u gorm.io/gorm
RUN go get -u github.com/gin-gonic/gin
RUN go get -u github.com/ilyakaznacheev/cleanenv

# Рабочая директория внутри контейнера

# Копируем go.mod и go.sum
# COPY go.sum ./
RUN go mod download

# Копируем всё остальное
# COPY . .

#  
RUN echo "go run cmd/server/main.go" > /root/.sh_history
# Устанавливаем переменную среды для увеличения размера истории
ENV HISTFILE=/root/.sh_history

CMD ["sh"]

# Собираем приложение
# RUN go build -o app ./cmd/server

# Указываем порт
EXPOSE 8080

# Команда запуска
# CMD ["./app"]
