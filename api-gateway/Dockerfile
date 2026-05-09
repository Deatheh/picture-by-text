# Этап сборки
FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY proto/ ./proto/
RUN cp -r ./proto/userpb ./userpb

COPY api-gateway/go.mod api-gateway/go.sum ./
RUN sed -i 's|../proto/userpb|./userpb|g' go.mod
RUN go mod download

COPY api-gateway/ ./api-gateway/

WORKDIR /build/api-gateway/cmd/app
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/api-gateway-app .

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем .env файл
COPY .env .env

# Копируем бинарник
COPY --from=builder /build/api-gateway-app .

EXPOSE 8080

CMD ["./api-gateway-app"]