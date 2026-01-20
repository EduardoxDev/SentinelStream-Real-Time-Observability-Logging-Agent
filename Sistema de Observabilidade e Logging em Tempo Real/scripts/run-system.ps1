$ErrorActionPreference = "Stop"
Write-Host "Observability System - Startup" -ForegroundColor Cyan
Write-Host "===============================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Verificando Redis..." -ForegroundColor Yellow
try {
    $redisTest = Test-NetConnection -ComputerName localhost -Port 6379 -WarningAction SilentlyContinue
    if ($redisTest.TcpTestSucceeded) {
        Write-Host "Redis: OK" -ForegroundColor Green
    } else {
        Write-Host "Redis nao esta rodando!" -ForegroundColor Red
        Write-Host "Execute: C:\observability\start-all.bat" -ForegroundColor Yellow
        exit 1
    }
} catch {
    Write-Host "Redis nao esta rodando!" -ForegroundColor Red
    exit 1
}
Write-Host "Verificando InfluxDB..." -ForegroundColor Yellow
try {
    $influxTest = Invoke-WebRequest -Uri "http://localhost:8086/health" -UseBasicParsing -TimeoutSec 3
    if ($influxTest.StatusCode -eq 200) {
        Write-Host "InfluxDB: OK" -ForegroundColor Green
    }
} catch {
    Write-Host "InfluxDB nao esta rodando!" -ForegroundColor Red
    Write-Host "Execute: C:\observability\start-all.bat" -ForegroundColor Yellow
    exit 1
}
Write-Host ""
Write-Host "Configuracao do InfluxDB" -ForegroundColor Cyan
Write-Host ""
$token = Read-Host "Cole o token do InfluxDB aqui (ou pressione Enter para usar o padrao)"
if ([string]::IsNullOrWhiteSpace($token)) {
    $token = "my-super-secret-token"
    Write-Host "Usando token padrao" -ForegroundColor Yellow
}
$env:INFLUXDB_URL = "http://localhost:8086"
$env:INFLUXDB_TOKEN = $token
$env:INFLUXDB_ORG = "observability"
$env:INFLUXDB_BUCKET = "metrics"
$env:REDIS_ADDR = "localhost:6379"
$env:PORT = "8080"
$env:CPU_THRESHOLD = "90.0"
$env:MEMORY_THRESHOLD = "85.0"
Write-Host ""
Write-Host "Variaveis de ambiente configuradas!" -ForegroundColor Green
Write-Host ""
$currentDir = Get-Location
Write-Host "Iniciando sistema..." -ForegroundColor Cyan
Write-Host ""
Write-Host "Abrindo 2 terminais:" -ForegroundColor Yellow
Write-Host "  Terminal 1: Agent (coleta metricas)" -ForegroundColor White
Write-Host "  Terminal 2: Server (API + WebSocket)" -ForegroundColor White
Write-Host ""
$agentCmd = "Write-Host 'Agent' -ForegroundColor Cyan; `$env:INFLUXDB_URL='$env:INFLUXDB_URL'; `$env:INFLUXDB_TOKEN='$token'; `$env:INFLUXDB_ORG='observability'; `$env:INFLUXDB_BUCKET='metrics'; `$env:REDIS_ADDR='localhost:6379'; Set-Location '$currentDir'; go run cmd/agent/main.go"
Start-Process powershell -ArgumentList "-NoExit", "-Command", $agentCmd
Start-Sleep -Seconds 2
$serverCmd = "Write-Host 'Server' -ForegroundColor Cyan; `$env:INFLUXDB_URL='$env:INFLUXDB_URL'; `$env:INFLUXDB_TOKEN='$token'; `$env:INFLUXDB_ORG='observability'; `$env:INFLUXDB_BUCKET='metrics'; `$env:REDIS_ADDR='localhost:6379'; `$env:PORT='8080'; Set-Location '$currentDir'; go run cmd/server/main.go"
Start-Process powershell -ArgumentList "-NoExit", "-Command", $serverCmd
Start-Sleep -Seconds 5
Write-Host ""
Write-Host "Sistema iniciado!" -ForegroundColor Green
Write-Host ""
Write-Host "Dashboards:" -ForegroundColor Cyan
Write-Host "  http://localhost:8080" -ForegroundColor White
Write-Host ""
Start-Process "http://localhost:8080"
Write-Host "Pronto!" -ForegroundColor Green