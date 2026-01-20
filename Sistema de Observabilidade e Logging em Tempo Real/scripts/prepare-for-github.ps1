$ErrorActionPreference = "Stop"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Prepare Repository for GitHub" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "[1/5] Removing build artifacts..." -ForegroundColor Yellow
if (Test-Path "bin") {
    Remove-Item -Path "bin" -Recurse -Force
    Write-Host "  Removed bin/" -ForegroundColor Green
}

Write-Host ""
Write-Host "[2/5] Removing temporary files..." -ForegroundColor Yellow
Get-ChildItem -Path . -Include *.exe,*.dll,*.so,*.dylib -Recurse | Remove-Item -Force
Write-Host "  Removed executables" -ForegroundColor Green

Write-Host ""
Write-Host "[3/5] Removing logs..." -ForegroundColor Yellow
if (Test-Path "logs") {
    Remove-Item -Path "logs" -Recurse -Force
    Write-Host "  Removed logs/" -ForegroundColor Green
}
Get-ChildItem -Path . -Include *.log -Recurse | Remove-Item -Force
Write-Host "  Removed log files" -ForegroundColor Green

Write-Host ""
Write-Host "[4/5] Removing cache..." -ForegroundColor Yellow
if (Test-Path ".cache") {
    Remove-Item -Path ".cache" -Recurse -Force
    Write-Host "  Removed .cache/" -ForegroundColor Green
}
if (Test-Path "tmp") {
    Remove-Item -Path "tmp" -Recurse -Force
    Write-Host "  Removed tmp/" -ForegroundColor Green
}

Write-Host ""
Write-Host "[5/5] Checking Git status..." -ForegroundColor Yellow

try {
    $gitStatus = git status --porcelain 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        $fileCount = ($gitStatus | Measure-Object).Count
        Write-Host "  Files to commit: $fileCount" -ForegroundColor White
        
        if ($fileCount -gt 0) {
            Write-Host ""
            Write-Host "Files ready for commit:" -ForegroundColor Cyan
            git status --short
        }
    }
} catch {
    Write-Host "  Git not initialized or not available" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "Repository Prepared!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "  1. git add ." -ForegroundColor White
Write-Host "  2. git commit -m 'Initial commit'" -ForegroundColor White
Write-Host "  3. git push origin main" -ForegroundColor White
Write-Host ""
Write-Host "If commit fails, try:" -ForegroundColor Yellow
Write-Host "  git config http.postBuffer 524288000" -ForegroundColor White
Write-Host ""
