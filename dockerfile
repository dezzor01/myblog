# === Этап сборки ===
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем статический бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /myblog cmd/myblog/main.go

# === Финальный образ (очень маленький!) ===
FROM alpine:3.20

# Устанавливаем только ca-certificates (для HTTPS)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем только бинарник и шаблоны
COPY --from=builder /myblog .
COPY --from=builder /app/internal/templates ./internal/templates

# Порт
EXPOSE 3000

# Запуск
CMD ["./myblog"]
