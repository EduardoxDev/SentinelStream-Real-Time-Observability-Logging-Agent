$ErrorActionPreference = "Stop"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Safe Code Cleaner" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "This script safely removes comments while preserving code structure." -ForegroundColor Yellow
Write-Host ""

$confirmation = Read-Host "Continue? (yes/no)"
if ($confirmation -ne "yes") {
    Write-Host "Aborted." -ForegroundColor Red
    exit 0
}

Write-Host ""

$stats = @{
    FilesProcessed = 0
    CommentsRemoved = 0
}

function Clean-GoFile {
    param($filePath)
    
    $lines = Get-Content $filePath
    $newLines = @()
    $inMultilineComment = $false
    
    foreach ($line in $lines) {
        $processedLine = $line
        
        if ($inMultilineComment) {
            if ($line -match '\*/') {
                $inMultilineComment = $false
                $processedLine = $line -replace '^.*?\*/', ''
                if ([string]::IsNullOrWhiteSpace($processedLine)) {
                    $stats.CommentsRemoved++
                    continue
                }
            } else {
                $stats.CommentsRemoved++
                continue
            }
        }
        
        if ($processedLine -match '/\*') {
            if ($processedLine -match '/\*.*?\*/') {
                $processedLine = $processedLine -replace '/\*.*?\*/', ''
            } else {
                $inMultilineComment = $true
                $processedLine = $processedLine -replace '/\*.*$', ''
            }
        }
        
        if ($processedLine -match '^\s*//') {
            $stats.CommentsRemoved++
            continue
        }
        
        if ($processedLine -match '\s+//\s+') {
            $processedLine = $processedLine -replace '\s+//.*$', ''
        }
        
        $newLines += $processedLine.TrimEnd()
    }
    
    $newLines -join "`n" | Set-Content $filePath -NoNewline
}

function Clean-JSFile {
    param($filePath)
    
    $lines = Get-Content $filePath
    $newLines = @()
    $inMultilineComment = $false
    
    foreach ($line in $lines) {
        $processedLine = $line
        
        if ($inMultilineComment) {
            if ($line -match '\*/') {
                $inMultilineComment = $false
                $processedLine = $line -replace '^.*?\*/', ''
                if ([string]::IsNullOrWhiteSpace($processedLine)) {
                    $stats.CommentsRemoved++
                    continue
                }
            } else {
                $stats.CommentsRemoved++
                continue
            }
        }
        
        if ($processedLine -match '/\*') {
            if ($processedLine -match '/\*.*?\*/') {
                $processedLine = $processedLine -replace '/\*.*?\*/', ''
            } else {
                $inMultilineComment = $true
                $processedLine = $processedLine -replace '/\*.*$', ''
            }
        }
        
        if ($processedLine -match '^\s*//') {
            $stats.CommentsRemoved++
            continue
        }
        
        $newLines += $processedLine.TrimEnd()
    }
    
    $newLines -join "`n" | Set-Content $filePath -NoNewline
}

function Clean-HTMLFile {
    param($filePath)
    
    $content = Get-Content $filePath -Raw
    $originalCount = ([regex]::Matches($content, '<!--')).Count
    $content = $content -replace '<!--.*?-->', ''
    $stats.CommentsRemoved += $originalCount
    $content | Set-Content $filePath -NoNewline
}

Write-Host "Processing Go files..." -ForegroundColor Cyan
Get-ChildItem -Path . -Include *.go -Recurse -File | Where-Object {
    $_.FullName -notmatch '\\vendor\\' -and 
    $_.FullName -notmatch '\\proto\\gen\\' -and
    $_.FullName -notmatch '\\bin\\'
} | ForEach-Object {
    Write-Host "  $($_.Name)" -ForegroundColor Gray
    Clean-GoFile $_.FullName
    $stats.FilesProcessed++
}

Write-Host "Processing JavaScript files..." -ForegroundColor Cyan
Get-ChildItem -Path . -Include *.js -Recurse -File | Where-Object {
    $_.FullName -notmatch '\\node_modules\\'
} | ForEach-Object {
    Write-Host "  $($_.Name)" -ForegroundColor Gray
    Clean-JSFile $_.FullName
    $stats.FilesProcessed++
}

Write-Host "Processing HTML files..." -ForegroundColor Cyan
Get-ChildItem -Path . -Include *.html -Recurse -File | ForEach-Object {
    Write-Host "  $($_.Name)" -ForegroundColor Gray
    Clean-HTMLFile $_.FullName
    $stats.FilesProcessed++
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "Completed!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Files processed: $($stats.FilesProcessed)" -ForegroundColor White
Write-Host "Comments removed: $($stats.CommentsRemoved)" -ForegroundColor White
Write-Host ""
Write-Host "Testing compilation..." -ForegroundColor Cyan

try {
    $agentBuild = go build -o bin/agent-test.exe ./cmd/agent 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  Agent: OK" -ForegroundColor Green
        Remove-Item bin/agent-test.exe -ErrorAction SilentlyContinue
    } else {
        Write-Host "  Agent: FAILED" -ForegroundColor Red
        Write-Host $agentBuild -ForegroundColor Red
    }
    
    $serverBuild = go build -o bin/server-test.exe ./cmd/server 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  Server: OK" -ForegroundColor Green
        Remove-Item bin/server-test.exe -ErrorAction SilentlyContinue
    } else {
        Write-Host "  Server: FAILED" -ForegroundColor Red
        Write-Host $serverBuild -ForegroundColor Red
    }
} catch {
    Write-Host "  Build test failed: $_" -ForegroundColor Red
}

Write-Host ""
Write-Host "Done! Your code is ready for GitHub." -ForegroundColor Green
Write-Host ""
