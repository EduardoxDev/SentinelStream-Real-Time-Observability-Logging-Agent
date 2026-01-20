<div align="center">

# ED Observability System

Real-time system monitoring platform with clean architecture

![Release](https://img.shields.io/badge/release-v1.0.0-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Go](https://img.shields.io/badge/go-1.21+-00ADD8.svg)

</div>

---

## Features

- 🎯 Real-time monitoring with WebSocket streaming
- 📊 Minimalist dashboard with 7 specialized views
- 🏗️ Clean Architecture (Hexagonal)
- ⚡ High performance Go backend
- 🔔 Smart alerts with configurable thresholds
- 📈 InfluxDB time series storage
- 🔄 Circuit breaker and retry policies
- 🌐 Windows-native support

---

## Quick Start

### Prerequisites

- Go 1.21+
- Redis 7.2+
- InfluxDB 2.7+

### Installation

```bash
git clone https://github.com/yourusername/observability-system.git
cd observability-system
go mod download
```

### Setup (Windows)

```powershell
.\scripts\setup-windows.ps1
```

Configure InfluxDB at http://localhost:8086:
- Organization: `observability`
- Bucket: `metrics`

### Run

```powershell
.\scripts\run-system.ps1
```

Access dashboard at http://localhost:8080

---

## Architecture

```
Presentation → Application → Domain → Infrastructure
    (UI)      (Use Cases)  (Entities)  (Adapters)
```

### Structure

```
cmd/          # Entry points
internal/     # Core logic
  ├─ domain/
  ├─ application/
  └─ infrastructure/
web/          # Frontend
scripts/      # Automation
```

---

## Configuration

```bash
INFLUXDB_TOKEN=your-token
CPU_THRESHOLD=90.0
MEMORY_THRESHOLD=85.0
```

---

## Development

Build:
```bash
go build ./cmd/agent
go build ./cmd/server
```

Test:
```bash
go test ./...
```

Clean code:
```powershell
.\scripts\safe-clean.ps1
```

---

## Deployment

**Docker:**
```bash
docker-compose up -d
```

**Kubernetes:**
```bash
kubectl apply -f k8s/
```

**AWS:**
```bash
cd terraform && terraform apply
```

---

## Performance

| Metric | Value |
|--------|-------|
| Agent RAM | ~15MB |
| Server RAM | ~25MB |
| CPU Usage | <2% |

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

---

## License

MIT License - see [LICENSE](LICENSE)

---

<div align="center">

**Built with Go, InfluxDB, and Redis**

Give a ⭐️ if this project helped you!

</div>
