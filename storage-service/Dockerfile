# Этап сборки
FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY proto/ ./proto/
RUN cp -r ./proto/storagepb ./storagepb

COPY storage-service/go.mod storage-service/go.sum ./
RUN sed -i 's|../proto/storagepb|./storagepb|g' go.mod
RUN go mod download

COPY storage-service/ ./storage-service/

# Сборка — явно указываем имя выходного файла
WORKDIR /build/storage-service/cmd/app
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/storage-service-app .

# Финальный образ
FROM alpine:latest

WORKDIR /app

COPY .env .env

# Копируем бинарник с правильным именем
COPY --from=builder /build/storage-service-app .

EXPOSE 50053

CMD ["./storage-service-app"]