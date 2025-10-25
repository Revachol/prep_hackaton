#!/bin/bash

echo "Generating Swagger documentation..."

# Генерируем Swagger документацию
swag init -g cmd/server/main.go -o docs/

if [ $? -eq 0 ]; then
    echo "✅ Swagger documentation generated successfully!"
    echo "📚 Docs will be available at: http://localhost:8080/swagger/index.html"
else
    echo "❌ Failed to generate Swagger documentation"
    exit 1
fi