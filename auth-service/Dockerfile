# Build stage
FROM golang:1.25-alpine AS builder

# Устанавливаем зависимости для компиляции и protoc
RUN apk add --no-cache git protoc protobuf-dev

# Устанавливаем swag для генерации документации
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Устанавливаем gRPC плагины и добавляем в PATH
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Добавляем GOPATH/bin в PATH для protoc плагинов
ENV PATH $PATH:/go/bin

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Генерируем Swagger документацию
RUN swag init -g cmd/server/main.go -o docs/

# Генерируем gRPC код (теперь с правильным PATH)
RUN mkdir -p gen && \
    protoc --go_out=gen --go_opt=paths=source_relative \
           --go-grpc_out=gen --go-grpc_opt=paths=source_relative \
           proto/auth.proto

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config
RUN adduser -D -s /bin/sh appuser
USER appuser
EXPOSE 8080 50051
CMD ["./main"]