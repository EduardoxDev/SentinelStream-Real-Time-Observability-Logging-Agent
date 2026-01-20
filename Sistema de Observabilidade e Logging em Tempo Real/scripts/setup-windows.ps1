Write-Host "Observability System - Windows Setup" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""
$baseDir = "C:\observability"
$redisDir = "$baseDir\redis"
$influxDir = "$baseDir\influxdb"
Write-Host "Criando diretorios..." -ForegroundColor Yellow
New-Item -ItemType Directory -Force -Path $baseDir | Out-Null
New-Item -ItemType Directory -Force -Path $redisDir | Out-Null
New-Item -ItemType Directory -Force -Path $influxDir | Out-Null
Write-Host "Baixando Redis..." -ForegroundColor Yellow
$redisUrl = "https://github.com/microsoftarchive/redis/releases/download/win-3.0.504/Redis-x64-3.0.504.zip"
$redisZip = "$baseDir\redis.zip"
if (-not (Test-Path "$redisDir\redis-server.exe")) {
    Invoke-WebRequest -Uri $redisUrl -OutFile $redisZip
    Expand-Archive -Path $redisZip -DestinationPath $redisDir -Force
    Remove-Item $redisZip
    Write-Host "Redis baixado!" -ForegroundColor Green
} else {
    Write-Host "Redis ja existe!" -ForegroundColor Green
}
Write-Host "Baixando InfluxDB..." -ForegroundColor Yellow
$influxUrl = "https://dl.influxdata.com/influxdb/releases/influxdb2-2.7.4-windows.zip"
$influxZip = "$baseDir\influxdb.zip"
if (-not (Test-Path "$influxDir\influxd.exe")) {
    Invoke-WebRequest -Uri $influxUrl -OutFile $influxZip
    Expand-Archive -Path $influxZip -DestinationPath $influxDir -Force
    Remove-Item $influxZip
    Write-Host "InfluxDB baixado!" -ForegroundColor Green
} else {
    Write-Host "InfluxDB ja existe!" -ForegroundColor Green
}
Write-Host "Criando scripts de inicializacao..." -ForegroundColor Yellow
$redisScript = @"
@echo off
cd /d $redisDir
echo Iniciando Redis...
start redis-server.exe redis.windows.conf
echo Redis iniciado na porta 6379
"@
$redisScript | Out-File -FilePath "$baseDir\start-redis.bat" -Encoding ASCII
$influxScript = @"
@echo off
cd /d $influxDir
echo Iniciando InfluxDB...
start influxd.exe
echo InfluxDB iniciado na porta 8086
echo Acesse: http://localhost:8086
"@
$influxScript | Out-File -FilePath "$baseDir\start-influxdb.bat" -Encoding ASCII
$startAllScript = @"
@echo off
echo ========================================
echo   Observability System - Starting
echo ========================================
echo.
echo [1/2] Iniciando Redis...
start /min cmd /c "$baseDir\start-redis.bat"
timeout /t 3 /nobreak >nul
echo [2/2] Iniciando InfluxDB...
start /min cmd /c "$baseDir\start-influxdb.bat"
timeout /t 5 /nobreak >nul
echo.
echo ========================================
echo   Servicos Iniciados!
echo ========================================
echo.
echo Redis:    localhost:6379
echo InfluxDB: http://localhost:8086
echo.
echo Configure o InfluxDB:
echo 1. Abra http://localhost:8086
echo 2. Username: admin
echo 3. Password: adminpass123
echo 4. Organization: observability
echo 5. Bucket: metrics
echo 6. Copie o token gerado
echo.
pause
"@
$startAllScript | Out-File -FilePath "$baseDir\start-all.bat" -Encoding ASCII
Write-Host ""
Write-Host "Instalacao concluida!" -ForegroundColor Green
Write-Host ""
Write-Host "Proximos passos:" -ForegroundColor Cyan
Write-Host "1. Execute: $baseDir\start-all.bat" -ForegroundColor White
Write-Host "2. Configure InfluxDB em http://localhost:8086" -ForegroundColor White
Write-Host "3. Copie o token gerado" -ForegroundColor White
Write-Host "4. Execute: .\scripts\run-system.ps1" -ForegroundColor White
Write-Host ""