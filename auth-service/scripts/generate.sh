#!/bin/bash

# Создаем директорию для сгенерированного кода
mkdir -p gen

# Генерируем код
protoc --go_out=gen --go_opt=paths=source_relative \
       --go-grpc_out=gen --go-grpc_opt=paths=source_relative \
       proto/auth.proto

echo "✅ gRPC code generated successfully!"