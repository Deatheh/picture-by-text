# Этап сборки
FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY proto/ ./proto/
RUN cp -r ./proto/userpb ./userpb

COPY user-service/go.mod user-service/go.sum ./
RUN sed -i 's|../proto/userpb|./userpb|g' go.mod
RUN go mod download

COPY user-service/ ./user-service/

# Сборка — явно указываем имя выходного файла
WORKDIR /build/user-service/cmd/app
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/user-service-app .

# Финальный образ
FROM alpine:latest

WORKDIR /app

COPY .env .env

# Копируем бинарник с правильным именем
COPY --from=builder /build/user-service-app .

EXPOSE 50051

CMD ["./user-service-app"]