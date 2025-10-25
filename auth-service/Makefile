.PHONY: build up down logs clean test swagger

# Сборка и запуск
build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

# Полный перезапуск
restart: down up

# Очистка
clean:
	docker-compose down -v
	docker system prune -f

# Генерация Swagger документации
swagger:
	swag init -g cmd/server/main.go -o docs/

# Генерация gRPC кода
grpc:
	mkdir -p gen
	protoc --go_out=gen --go_opt=paths=source_relative \
	       --go-grpc_out=gen --go-grpc_opt=paths=source_relative \
	       proto/auth.proto

# Тестирование
test:
	go test ./...

# Запуск в development режиме
dev: swagger grpc
	go run cmd/server/main.go

# Просмотр логов конкретного сервиса
logs-auth:
	docker-compose logs -f auth-service

logs-db:
	docker-compose logs -f postgres

logs-redis:
	docker-compose logs -f redis

# Проверка статуса сервисов
status:
	docker-compose ps

# Bash в контейнер сервиса
bash:
	docker-compose exec auth-service /bin/sh