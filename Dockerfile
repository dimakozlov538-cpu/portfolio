# Этап 1: Сборка
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Сначала копируем только файлы зависимостей
COPY go.mod ./
COPY go.sum ./

# Скачиваем зависимости
RUN go mod download

# Теперь копируем весь код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o portfolio .

# Этап 2: Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем бинарник из builder
COPY --from=builder /app/portfolio .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Устанавливаем часовой пояс
ENV TZ=Europe/Moscow

EXPOSE 8080

CMD ["./portfolio"]