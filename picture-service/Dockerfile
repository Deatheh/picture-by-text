# Этап сборки
FROM golang:1.26-alpine AS builder

WORKDIR /build

# Копируем оба proto (picturepb для API, storagepb для вызовов)
COPY proto/ ./proto/
RUN cp -r ./proto/picturepb ./picturepb

# Копируем go.mod и go.sum
COPY picture-service/go.mod picture-service/go.sum ./

# Заменяем оба replace
RUN sed -i 's|../proto/picturepb|./picturepb|g' go.mod

RUN go mod download

# Копируем исходники
COPY picture-service/ ./picture-service/

# Сборка
WORKDIR /build/picture-service/cmd/app
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/picture-service-app .

# Финальный образ
FROM alpine:latest

WORKDIR /app

COPY .env .env

# Копируем бинарник
COPY --from=builder /build/picture-service-app .

EXPOSE 50052

CMD ["./picture-service-app"]