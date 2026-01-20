<div align="center">

# 🔍 Observability System

### Real-time Docker Container Monitoring & Alerting Platform

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker)](https://docker.com)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-Supported-326CE5?style=for-the-badge&logo=kubernetes)](https://kubernetes.io)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

[Features](#-features) • [Architecture](#-architecture) • [Quick Start](#-quick-start) • [Documentation](#-documentation) • [Benchmarks](#-performance)

![Observability System](docs/images/system-overview.png)

</div>

---

## 💡 About

A production-ready observability platform built from scratch to monitor Docker containers with real-time metrics collection, intelligent alerting, and beautiful dashboards. Designed with enterprise patterns like Clean Architecture, Circuit Breaker, and full Kubernetes support.

## ✨ Features

- 🚀 **Real-time Monitoring** - Live metrics via WebSockets with sub-second latency
- 📊 **Historical Analytics** - Time-series data with interactive Chart.js dashboards
- 🔔 **Smart Alerts** - Multi-channel notifications (Slack, Discord, Email) with cooldown
- 🛡️ **Resilient** - Circuit Breaker & Retry patterns for fault tolerance
- ⚡ **High Performance** - 4,260 ops/sec with minimal resource footprint
- 🏗️ **Clean Architecture** - Hexagonal design with clear separation of concerns
- ☸️ **Cloud Native** - Full Kubernetes support with Helm charts
- 🔐 **Secure** - JWT authentication with role-based access control
- 📈 **Prometheus Ready** - Native metrics export for Prometheus scraping
- 🌐 **gRPC Enabled** - High-performance inter-service communication

## 🛠️ Tech Stack

| Category | Technology |
|----------|-----------|
| **Language** | Go 1.21+ |
| **Container Runtime** | Docker API |
| **Time-Series DB** | InfluxDB 2.x |
| **Cache & Messaging** | Redis 7.x |
| **Real-time Comms** | WebSockets, gRPC |
| **Frontend** | Vanilla JS, Chart.js |
| **Infrastructure** | Terraform (AWS), Kubernetes, Helm |
| **Monitoring** | Prometheus, CloudWatch |

## �️ Arcuhitecture

Built with **Clean Architecture** principles for maintainability and testability:

```
┌─────────────────────────────────────────────────────────┐
│                    Presentation Layer                    │
│              (HTTP, WebSocket, gRPC APIs)               │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│                  Application Layer                       │
│         (Use Cases: Collect, Alert, Analyze)            │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│                    Domain Layer                          │
│        (Entities, Business Rules, Interfaces)           │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│                Infrastructure Layer                      │
│   (Docker, InfluxDB, Redis, Notifiers, Resilience)     │
└─────────────────────────────────────────────────────────┘
```

**Key Patterns:**
- 🎯 SOLID Principles
- 🔌 Ports & Adapters (Hexagonal)
- 🛡️ Circuit Breaker & Retry
- 📦 Dependency Injection

## � Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+
- Make (optional)

### 1. Clone & Setup
```bash
git clone https://github.com/yourusername/observability-system.git
cd observability-system
go mod download
```

### 2. Start Infrastructure
```bash
docker-compose up -d
```

### 3. Generate gRPC Code
```bash
# Windows
.\scripts\generate_proto.ps1

# Linux/Mac
chmod +x scripts/generate_proto.sh && ./scripts/generate_proto.sh
```

### 4. Run the System
```bash
# Terminal 1 - Agent (collects metrics)
go run cmd/agent/main.go

# Terminal 2 - Server (API & WebSocket)
go run cmd/server/main.go
```

### 5. Access Dashboards
- 🌐 **Real-time Dashboard**: http://localhost:8080
- 📊 **Historical Charts**: http://localhost:8080/dashboard.html
- 📈 **Prometheus Metrics**: http://localhost:9090/metrics

That's it! You should see metrics flowing in real-time.

## 🏗️ Arquitetura

```
┌─────────────┐     ┌──────────────┐     ┌──────────────┐
│   Docker    │────▶│    Agent     │────▶│   InfluxDB   │
│  Containers │     │  (Collector) │     │  (Storage)   │
└─────────────┘     └──────────────┘     └──────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │    Redis     │
                    │  (Alerting)  │
                    └──────────────┘
                           │
                           ▼
                    ┌──────────────┐     ┌──────────────┐
                    │    Server    │────▶│  Dashboard   │
                    │  (WebSocket) │     │   (Web UI)   │
                    └──────────────┘     └──────────────┘
```

## � API tReference

### REST Endpoints
```http
GET  /api/containers              # List running containers
GET  /api/metrics?container_id=x  # Historical metrics
POST /api/auth/login              # JWT authentication
GET  /health                      # Health check
```

### WebSocket
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onmessage = (event) => {
  const metrics = JSON.parse(event.data);
  console.log(metrics);
};
```

### gRPC
```protobuf
service MetricsService {
  rpc StreamMetrics(stream MetricData) returns (stream MetricResponse);
  rpc SubscribeToMetrics(SubscriptionRequest) returns (stream MetricData);
}
```

## ⚙️ Configuration

Create a `.env` file or set environment variables:

```bash
# InfluxDB Configuration
INFLUXDB_URL=http://localhost:8086
INFLUXDB_TOKEN=your-secret-token
INFLUXDB_ORG=observability
INFLUXDB_BUCKET=metrics

# Redis Configuration
REDIS_ADDR=localhost:6379

# Server Configuration
PORT=8080
GRPC_PORT=50051

# Alert Thresholds
CPU_THRESHOLD=90.0
MEMORY_THRESHOLD=85.0

# JWT Authentication
JWT_SECRET=your-jwt-secret
JWT_DURATION=24h

# Notifiers (optional)
SLACK_WEBHOOK_URL=https://hooks.slack.com/...
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/...
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
```

## 🚨 Sistema de Alertas

O sistema monitora automaticamente:
- **CPU > 90%** - Alerta crítico
- **Memória > 85%** - Alerta de memória

Alertas possuem cooldown de 5 minutos para evitar spam.

## 📈 Métricas Coletadas

- CPU Usage (%)
- Memory Usage (bytes e %)
- Network RX/TX (bytes)
- Timestamp de coleta

## 🧪 Testando

Para gerar carga em um container:
```bash
docker run -d --name stress-test progrium/stress --cpu 2 --vm 1 --vm-bytes 128M
```

## 🔧 Desenvolvimento

### Gerar código gRPC
```bash
protoc --go_out=. --go-grpc_out=. proto/metrics.proto
```

### Build
```bash
go build -o bin/agent cmd/agent/main.go
go build -o bin/server cmd/server/main.go
```

## �️ Róesiliência e Padrões

### Circuit Breaker

Protege o sistema contra falhas em cascata:

```go
// Configuração padrão
MaxFailures: 5        // Abre após 5 falhas
Timeout: 30s          // Tenta reconectar após 30s
States: CLOSED → OPEN → HALF_OPEN → CLOSED
```

**Comportamento:**
- **CLOSED**: Operação normal
- **OPEN**: Fail-fast, não tenta operação (protege recursos)
- **HALF_OPEN**: Testa se serviço voltou

### Retry Policy

Retry automático com backoff exponencial:

```go
MaxAttempts: 3        // Máximo de tentativas
InitialDelay: 1s      // Delay inicial
Backoff: 2.0          // Multiplica delay a cada tentativa
```

**Exemplo:**
- Tentativa 1: Falha → Aguarda 1s
- Tentativa 2: Falha → Aguarda 2s
- Tentativa 3: Sucesso ✅

### O que acontece se o banco cair?

1. **InfluxDB indisponível:**
   - Circuit breaker abre após 5 falhas
   - Requisições falham rapidamente (fail-fast)
   - Sistema continua coletando métricas
   - Após 30s, tenta reconectar automaticamente
   - Quando volta, retoma operação normal

2. **Redis indisponível:**
   - Retry policy tenta 3x com backoff
   - Alertas não são enviados (graceful degradation)
   - Coleta de métricas continua funcionando
   - Logs indicam problema com alertas

## 📊 Performance e Benchmarks

Veja o relatório completo em [BENCHMARKS.md](BENCHMARKS.md)

### Resumo dos Resultados

| Métrica | Valor |
|---------|-------|
| Throughput (InfluxDB) | 4,260 ops/s |
| Latência P95 | 3.2ms |
| Memória (10 containers) | 8.2 MB |
| CPU (100 containers) | 8.2% |
| Taxa de Sucesso | 99.97% |

**Executar benchmarks:**
```bash
chmod +x scripts/run_benchmarks.sh
./scripts/run_benchmarks.sh
```

## ☁️ Deploy na AWS (Terraform)

### Infraestrutura Provisionada

- **VPC** com subnets públicas e privadas
- **ECS Fargate** para Agent e Server
- **Application Load Balancer** para distribuição de carga
- **ElastiCache Redis** para alertas
- **Secrets Manager** para credenciais
- **CloudWatch** para logs e métricas
- **IAM Roles** com least privilege

### Deploy

```bash
cd terraform

# Inicializar Terraform
terraform init

# Planejar mudanças
terraform plan -var="ecr_repository_url=YOUR_ECR_URL" \
               -var="influxdb_token=YOUR_TOKEN"

# Aplicar infraestrutura
terraform apply

# Obter outputs
terraform output alb_dns_name
```

### Variáveis Necessárias

```hcl
aws_region         = "us-east-1"
environment        = "production"
ecr_repository_url = "123456789.dkr.ecr.us-east-1.amazonaws.com/observability"
influxdb_token     = "your-secret-token"
```

### Custos Estimados (AWS)

- ECS Fargate (2 tasks): ~$30/mês
- ElastiCache (t3.micro): ~$15/mês
- ALB: ~$20/mês
- **Total**: ~$65/mês

## 🐳 Build e Deploy com Docker

```bash
# Build das imagens
docker build -f Dockerfile.agent -t observability-agent:latest .
docker build -f Dockerfile.server -t observability-server:latest .

# Push para ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin YOUR_ECR_URL
docker tag observability-agent:latest YOUR_ECR_URL:agent-latest
docker tag observability-server:latest YOUR_ECR_URL:server-latest
docker push YOUR_ECR_URL:agent-latest
docker push YOUR_ECR_URL:server-latest
```

## 📄 Licença

MIT
