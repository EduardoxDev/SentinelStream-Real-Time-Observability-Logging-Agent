.PHONY: help setup run-infra run-agent run-server stop clean

help:
	@echo "🔍 Observability System - Comandos disponíveis:"
	@echo "  make setup       - Instala dependências"
	@echo "  make run-infra   - Inicia Redis e InfluxDB"
	@echo "  make run-agent   - Executa o agent de coleta"
	@echo "  make run-server  - Executa o servidor web"
	@echo "  make stop        - Para todos os serviços"
	@echo "  make clean       - Remove volumes e dados"

setup:
	go mod download
	go mod tidy

run-infra:
	docker-compose up -d
	@echo "✅ Infraestrutura iniciada!"
	@echo "📊 InfluxDB: http://localhost:8086"
	@echo "🔴 Redis: localhost:6379"

run-agent:
	go run cmd/agent/main.go

run-server:
	go run cmd/server/main.go

stop:
	docker-compose down

clean:
	docker-compose down -v
	rm -rf bin/

build:
	mkdir -p bin
	go build -o bin/agent cmd/agent/main.go
	go build -o bin/server cmd/server/main.go
