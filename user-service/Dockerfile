# Этап сборки
FROM golang:1.26-alpine AS builder

WORKDIR /build

# Копируем оба proto (userpb для API, storagepb для вызовов)
COPY proto/ ./proto/
RUN cp -r ./proto/userpb ./userpb

# Копируем go.mod и go.sum
COPY user-service/go.mod user-service/go.sum ./

# Заменяем оба replace
RUN sed -i 's|../proto/userpb|./userpb|g' go.mod

RUN go mod download

# Копируем исходники
COPY user-service/ ./user-service/

# Сборка
WORKDIR /build/user-service/cmd/app
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/user-service-app .

# Финальный образ
FROM alpine:latest

WORKDIR /app

COPY .env .env

# Копируем бинарник
COPY --from=builder /build/user-service-app .

EXPOSE 50051

CMD ["./user-service-app"]