$ErrorActionPreference = "Stop"
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Clean Code Script" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "This script will:" -ForegroundColor Yellow
Write-Host "  1. Remove all comments" -ForegroundColor White
Write-Host "  2. Remove excessive blank lines" -ForegroundColor White
Write-Host "  3. Trim trailing whitespace" -ForegroundColor White
Write-Host ""
Write-Host "WARNING: Make sure you have a backup!" -ForegroundColor Red
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
    LinesRemoved = 0
}
function Clean-GoFile {
    param($filePath)
    $lines = Get-Content $filePath
    $newLines = @()
    $blankCount = 0
    foreach ($line in $lines) {
        $trimmed = $line.TrimEnd()
        if ($trimmed -match '^\s*//') {
            $stats.CommentsRemoved++
            continue
        }
        $trimmed = $trimmed -replace '\s*//.*$', ''
        if ([string]::IsNullOrWhiteSpace($trimmed)) {
            $blankCount++
            if ($blankCount -le 1) {
                $newLines += ""
            } else {
                $stats.LinesRemoved++
            }
        } else {
            $blankCount = 0
            $newLines += $trimmed
        }
    }
    while ($newLines.Count -gt 0 -and [string]::IsNullOrWhiteSpace($newLines[-1])) {
        $newLines = $newLines[0..($newLines.Count - 2)]
        $stats.LinesRemoved++
    }
    $newLines -join "`n" | Set-Content $filePath -NoNewline
}
function Clean-JSFile {
    param($filePath)
    $content = Get-Content $filePath -Raw
    $content = $content -replace '(?m)^\s*//.*$', ''
    $content = $content -replace '(?s)/\*.*?\*/', ''
    $lines = $content -split "`n"
    $newLines = @()
    $blankCount = 0
    foreach ($line in $lines) {
        $trimmed = $line.TrimEnd()
        if ([string]::IsNullOrWhiteSpace($trimmed)) {
            $blankCount++
            if ($blankCount -le 1) {
                $newLines += ""
            } else {
                $stats.LinesRemoved++
            }
        } else {
            $blankCount = 0
            $newLines += $trimmed
        }
    }
    while ($newLines.Count -gt 0 -and [string]::IsNullOrWhiteSpace($newLines[-1])) {
        $newLines = $newLines[0..($newLines.Count - 2)]
        $stats.LinesRemoved++
    }
    $newLines -join "`n" | Set-Content $filePath -NoNewline
}
function Clean-HTMLFile {
    param($filePath)
    $content = Get-Content $filePath -Raw
    $content = $content -replace '<!--.*?-->', ''
    $content | Set-Content $filePath -NoNewline
}
function Clean-PSFile {
    param($filePath)
    $lines = Get-Content $filePath
    $newLines = @()
    $blankCount = 0
    foreach ($line in $lines) {
        $trimmed = $line.TrimEnd()
        if ($trimmed -match '^\s*#') {
            $stats.CommentsRemoved++
            continue
        }
        if ([string]::IsNullOrWhiteSpace($trimmed)) {
            $blankCount++
            if ($blankCount -le 1) {
                $newLines += ""
            } else {
                $stats.LinesRemoved++
            }
        } else {
            $blankCount = 0
            $newLines += $trimmed
        }
    }
    $newLines -join "`n" | Set-Content $filePath -NoNewline
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
Write-Host "Processing PowerShell files..." -ForegroundColor Cyan
Get-ChildItem -Path . -Include *.ps1 -Recurse -File | Where-Object {
    $_.Name -ne "clean-code.ps1" -and $_.Name -ne "remove-comments.ps1"
} | ForEach-Object {
    Write-Host "  $($_.Name)" -ForegroundColor Gray
    Clean-PSFile $_.FullName
    $stats.FilesProcessed++
}
Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "Completed!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Files processed: $($stats.FilesProcessed)" -ForegroundColor White
Write-Host "Comments removed: $($stats.CommentsRemoved)" -ForegroundColor White
Write-Host "Blank lines removed: $($stats.LinesRemoved)" -ForegroundColor White
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "  1. Run: go build ./..." -ForegroundColor White
Write-Host "  2. Run: go test ./..." -ForegroundColor White
Write-Host "  3. Review changes: git diff" -ForegroundColor White
Write-Host ""
