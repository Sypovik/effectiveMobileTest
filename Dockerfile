FROM golang:1.24-alpine

WORKDIR /app
COPY go.mod ./


RUN go mod download

# добавляем команду запуска в историю
RUN echo "go run cmd/server/main.go" > /root/.sh_history
ENV HISTFILE=/root/.sh_history

EXPOSE 8080
