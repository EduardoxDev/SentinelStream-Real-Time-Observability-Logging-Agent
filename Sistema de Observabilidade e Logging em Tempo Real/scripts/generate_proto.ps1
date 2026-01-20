Write-Host "🔧 Generating gRPC code from proto files..." -ForegroundColor Cyan
New-Item -ItemType Directory -Force -Path "proto\gen" | Out-Null
protoc --go_out=proto\gen --go_opt=paths=source_relative `
       --go-grpc_out=proto\gen --go-grpc_opt=paths=source_relative `
       proto\metrics.proto
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ gRPC code generated successfully!" -ForegroundColor Green
    Write-Host "Files created in proto\gen\" -ForegroundColor Green
} else {
    Write-Host "❌ Failed to generate gRPC code" -ForegroundColor Red
    Write-Host "Make sure protoc, protoc-gen-go, and protoc-gen-go-grpc are installed" -ForegroundColor Yellow
}