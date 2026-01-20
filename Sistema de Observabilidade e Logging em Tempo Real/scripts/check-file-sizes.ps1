$ErrorActionPreference = "Stop"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "File Size Checker for GitHub" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$maxSize = 50MB
$warningSize = 10MB
$largeFiles = @()
$warningFiles = @()

Write-Host "Scanning files..." -ForegroundColor Yellow
Write-Host ""

Get-ChildItem -Path . -Recurse -File | Where-Object {
    $_.FullName -notmatch '\\\.git\\' -and
    $_.FullName -notmatch '\\bin\\' -and
    $_.FullName -notmatch '\\node_modules\\' -and
    $_.FullName -notmatch '\\vendor\\' -and
    $_.FullName -notmatch '\\C:\\observability\\'
} | ForEach-Object {
    $size = $_.Length
    
    if ($size -gt $maxSize) {
        $largeFiles += [PSCustomObject]@{
            Path = $_.FullName.Replace((Get-Location).Path + '\', '')
            Size = [math]::Round($size / 1MB, 2)
            SizeStr = "{0:N2} MB" -f ($size / 1MB)
        }
    }
    elseif ($size -gt $warningSize) {
        $warningFiles += [PSCustomObject]@{
            Path = $_.FullName.Replace((Get-Location).Path + '\', '')
            Size = [math]::Round($size / 1MB, 2)
            SizeStr = "{0:N2} MB" -f ($size / 1MB)
        }
    }
}

if ($largeFiles.Count -gt 0) {
    Write-Host "CRITICAL - Files too large for GitHub (>50MB):" -ForegroundColor Red
    Write-Host ""
    $largeFiles | Sort-Object -Property Size -Descending | ForEach-Object {
        Write-Host "  $($_.SizeStr.PadLeft(10)) - $($_.Path)" -ForegroundColor Red
    }
    Write-Host ""
}

if ($warningFiles.Count -gt 0) {
    Write-Host "WARNING - Large files (>10MB):" -ForegroundColor Yellow
    Write-Host ""
    $warningFiles | Sort-Object -Property Size -Descending | ForEach-Object {
        Write-Host "  $($_.SizeStr.PadLeft(10)) - $($_.Path)" -ForegroundColor Yellow
    }
    Write-Host ""
}

if ($largeFiles.Count -eq 0 -and $warningFiles.Count -eq 0) {
    Write-Host "All files are within acceptable size limits!" -ForegroundColor Green
    Write-Host ""
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Summary" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Critical files (>50MB): $($largeFiles.Count)" -ForegroundColor $(if ($largeFiles.Count -gt 0) { "Red" } else { "Green" })
Write-Host "Warning files (>10MB): $($warningFiles.Count)" -ForegroundColor $(if ($warningFiles.Count -gt 0) { "Yellow" } else { "Green" })
Write-Host ""

if ($largeFiles.Count -gt 0) {
    Write-Host "ACTIONS REQUIRED:" -ForegroundColor Red
    Write-Host "1. Add large files to .gitignore" -ForegroundColor White
    Write-Host "2. Use Git LFS for binary files" -ForegroundColor White
    Write-Host "3. Split large files into smaller chunks" -ForegroundColor White
    Write-Host "4. Move large files to external storage" -ForegroundColor White
    Write-Host ""
}
