$ErrorActionPreference = "Stop"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Remove Comments Script" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "WARNING: This will remove ALL comments from your code!" -ForegroundColor Yellow
Write-Host "Make sure you have a backup or committed your changes to git." -ForegroundColor Yellow
Write-Host ""

$confirmation = Read-Host "Do you want to continue? (yes/no)"
if ($confirmation -ne "yes") {
    Write-Host "Aborted." -ForegroundColor Red
    exit 0
}

Write-Host ""
Write-Host "Processing files..." -ForegroundColor Cyan
Write-Host ""

$filesProcessed = 0
$commentsRemoved = 0

function Remove-GoComments {
    param($filePath)
    
    $content = Get-Content $filePath -Raw
    $originalLength = $content.Length
    
    $content = $content -replace '(?m)^\s*//.*$', ''
    $content = $content -replace '(?s)/\*.*?\*/', ''
    $content = $content -replace '(?m)^\s*$\n', ''
    
    $content | Set-Content $filePath -NoNewline
    
    $newLength = (Get-Content $filePath -Raw).Length
    return $originalLength - $newLength
}

function Remove-JSComments {
    param($filePath)
    
    $content = Get-Content $filePath -Raw
    $originalLength = $content.Length
    
    $content = $content -replace '(?m)^\s*//.*$', ''
    $content = $content -replace '(?s)/\*.*?\*/', ''
    $content = $content -replace '(?m)^\s*$\n', ''
    
    $content | Set-Content $filePath -NoNewline
    
    $newLength = (Get-Content $filePath -Raw).Length
    return $originalLength - $newLength
}

function Remove-HTMLComments {
    param($filePath)
    
    $content = Get-Content $filePath -Raw
    $originalLength = $content.Length
    
    $content = $content -replace '<!--.*?-->', ''
    
    $content | Set-Content $filePath -NoNewline
    
    $newLength = (Get-Content $filePath -Raw).Length
    return $originalLength - $newLength
}

function Remove-PSComments {
    param($filePath)
    
    $content = Get-Content $filePath -Raw
    $originalLength = $content.Length
    
    $content = $content -replace '(?m)^\s*#.*$', ''
    $content = $content -replace '(?s)<#.*?#>', ''
    $content = $content -replace '(?m)^\s*$\n', ''
    
    $content | Set-Content $filePath -NoNewline
    
    $newLength = (Get-Content $filePath -Raw).Length
    return $originalLength - $newLength
}

Write-Host "[1/4] Processing Go files..." -ForegroundColor Yellow
Get-ChildItem -Path . -Include *.go -Recurse -File | Where-Object {
    $_.FullName -notmatch '\\vendor\\' -and 
    $_.FullName -notmatch '\\node_modules\\' -and
    $_.FullName -notmatch '\\bin\\'
} | ForEach-Object {
    Write-Host "  Processing: $($_.Name)" -ForegroundColor Gray
    $removed = Remove-GoComments $_.FullName
    $commentsRemoved += $removed
    $filesProcessed++
}

Write-Host "[2/4] Processing JavaScript files..." -ForegroundColor Yellow
Get-ChildItem -Path . -Include *.js -Recurse -File | Where-Object {
    $_.FullName -notmatch '\\node_modules\\' -and
    $_.FullName -notmatch '\\dist\\'
} | ForEach-Object {
    Write-Host "  Processing: $($_.Name)" -ForegroundColor Gray
    $removed = Remove-JSComments $_.FullName
    $commentsRemoved += $removed
    $filesProcessed++
}

Write-Host "[3/4] Processing HTML files..." -ForegroundColor Yellow
Get-ChildItem -Path . -Include *.html -Recurse -File | ForEach-Object {
    Write-Host "  Processing: $($_.Name)" -ForegroundColor Gray
    $removed = Remove-HTMLComments $_.FullName
    $commentsRemoved += $removed
    $filesProcessed++
}

Write-Host "[4/4] Processing PowerShell files..." -ForegroundColor Yellow
Get-ChildItem -Path . -Include *.ps1 -Recurse -File | Where-Object {
    $_.Name -ne "remove-comments.ps1"
} | ForEach-Object {
    Write-Host "  Processing: $($_.Name)" -ForegroundColor Gray
    $removed = Remove-PSComments $_.FullName
    $commentsRemoved += $removed
    $filesProcessed++
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "Completed!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Files processed: $filesProcessed" -ForegroundColor White
Write-Host "Approximate characters removed: $commentsRemoved" -ForegroundColor White
Write-Host ""
Write-Host "IMPORTANT: Review your code before committing!" -ForegroundColor Yellow
Write-Host "Run 'go build' to ensure everything still compiles." -ForegroundColor Yellow
Write-Host ""
