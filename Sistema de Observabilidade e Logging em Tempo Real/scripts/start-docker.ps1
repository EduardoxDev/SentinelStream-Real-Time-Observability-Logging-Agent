Write-Host "🐳 Starting Docker Infrastructure..." -ForegroundColor Cyan
$dockerRunning = docker ps 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Docker não está rodando!" -ForegroundColor Red
    Write-Host "Iniciando Docker Desktop..." -ForegroundColor Yellow
    Start-Process "C:\Program Files\Docker\Docker\Docker Desktop.exe"
    Write-Host "Aguardando Docker iniciar (30s)..." -ForegroundColor Yellow
    Start-Sleep -Seconds 30
    $dockerRunning = docker ps 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Falha ao iniciar Docker Desktop" -ForegroundColor Red
        Write-Host "Por favor, inicie manualmente e execute este script novamente" -ForegroundColor Yellow
        exit 1
    }
}
Write-Host "✅ Docker está rodando!" -ForegroundColor Green
Write-Host "`n🧹 Limpando containers antigos..." -ForegroundColor Cyan
docker stop observability-redis observability-influxdb observability-demo-app 2>$null
docker rm observability-redis observability-influxdb observability-demo-app 2>$null
Write-Host "`n📥 Baixando imagens Docker..." -ForegroundColor Cyan
docker pull redis:7-alpine
docker pull influxdb:2.7
docker pull nginx:alpine
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Erro ao baixar imagens" -ForegroundColor Red
    Write-Host "Verifique sua conexão com a internet" -ForegroundColor Yellow
    exit 1
}
Write-Host "`n🔴 Iniciando Redis..." -ForegroundColor Cyan
docker run -d --name observability-redis `
    -p 6379:6379 `
    --restart unless-stopped `
    redis:7-alpine
Start-Sleep -Seconds 2
Write-Host "📊 Iniciando InfluxDB..." -ForegroundColor Cyan
docker run -d --name observability-influxdb `
    -p 8086:8086 `
    -e DOCKER_INFLUXDB_INIT_MODE=setup `
    -e DOCKER_INFLUXDB_INIT_USERNAME=admin `
    -e DOCKER_INFLUXDB_INIT_PASSWORD=adminpass123 `
    -e DOCKER_INFLUXDB_INIT_ORG=observability `
    -e DOCKER_INFLUXDB_INIT_BUCKET=metrics `
    -e DOCKER_INFLUXDB_INIT_ADMIN_TOKEN=my-super-secret-token `
    --restart unless-stopped `
    influxdb:2.7
Start-Sleep -Seconds 5
Write-Host "🌐 Iniciando Demo App (nginx)..." -ForegroundColor Cyan
docker run -d --name observability-demo-app `
    --restart unless-stopped `
    nginx:alpine
Write-Host "`n✅ Verificando containers..." -ForegroundColor Cyan
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
Write-Host "`n🔍 Testando conectividade..." -ForegroundColor Cyan
$redisTest = docker exec observability-redis redis-cli ping 2>&1
if ($redisTest -eq "PONG") {
    Write-Host "✅ Redis: OK" -ForegroundColor Green
} else {
    Write-Host "❌ Redis: FALHOU" -ForegroundColor Red
}
Write-Host "Aguardando InfluxDB inicializar (10s)..." -ForegroundColor Yellow
Start-Sleep -Seconds 10
try {
    $influxTest = Invoke-WebRequest -Uri "http://localhost:8086/health" -UseBasicParsing -TimeoutSec 5
    if ($influxTest.StatusCode -eq 200) {
        Write-Host "✅ InfluxDB: OK" -ForegroundColor Green
    }
} catch {
    Write-Host "⚠️  InfluxDB: Ainda inicializando (normal)" -ForegroundColor Yellow
}
Write-Host "`n🎉 Infraestrutura iniciada com sucesso!" -ForegroundColor Green
Write-Host "`nEndpoints disponíveis:" -ForegroundColor Cyan
Write-Host "  Redis:    localhost:6379" -ForegroundColor White
Write-Host "  InfluxDB: http://localhost:8086" -ForegroundColor White
Write-Host "  InfluxDB UI: http://localhost:8086 (admin/adminpass123)" -ForegroundColor White
Write-Host "`nPróximos passos:" -ForegroundColor Cyan
Write-Host "  1. go run cmd/agent/main.go" -ForegroundColor White
Write-Host "  2. go run cmd/server/main.go" -ForegroundColor White
Write-Host "  3. Abrir http://localhost:8080" -ForegroundColor White