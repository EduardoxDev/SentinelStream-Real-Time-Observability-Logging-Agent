$ErrorActionPreference = "Stop"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Commit Repository in Chunks" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$commitMessage = Read-Host "Enter commit message"
if ([string]::IsNullOrWhiteSpace($commitMessage)) {
    $commitMessage = "Update codebase"
}

Write-Host ""
Write-Host "Preparing repository..." -ForegroundColor Yellow
.\scripts\prepare-for-github.ps1

Write-Host ""
Write-Host "Committing in chunks..." -ForegroundColor Cyan
Write-Host ""

$chunks = @(
    @{
        Name = "Core Go files"
        Patterns = @("*.go", "go.mod", "go.sum")
    },
    @{
        Name = "Web files"
        Patterns = @("web/*")
    },
    @{
        Name = "Scripts"
        Patterns = @("scripts/*")
    },
    @{
        Name = "Configuration"
        Patterns = @("*.yml", "*.yaml", "*.json", "*.toml", ".gitignore", ".env.example")
    },
    @{
        Name = "Documentation"
        Patterns = @("*.md", "LICENSE")
    },
    @{
        Name = "Infrastructure"
        Patterns = @("terraform/*", "k8s/*", "helm/*", "Dockerfile*", "docker-compose.yml")
    },
    @{
        Name = "Proto files"
        Patterns = @("proto/*")
    }
)

$chunkNumber = 1
foreach ($chunk in $chunks) {
    Write-Host "[$chunkNumber/$($chunks.Count)] Committing: $($chunk.Name)" -ForegroundColor Yellow
    
    foreach ($pattern in $chunk.Patterns) {
        git add $pattern 2>$null
    }
    
    $status = git status --porcelain
    if ($status) {
        git commit -m "$commitMessage - $($chunk.Name)" 2>&1 | Out-Null
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  Committed successfully" -ForegroundColor Green
        } else {
            Write-Host "  Nothing to commit or error occurred" -ForegroundColor Gray
        }
    } else {
        Write-Host "  No changes to commit" -ForegroundColor Gray
    }
    
    $chunkNumber++
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "All chunks committed!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Push to GitHub:" -ForegroundColor Cyan
Write-Host "  git push origin main" -ForegroundColor White
Write-Host ""
Write-Host "Or push with increased buffer:" -ForegroundColor Yellow
Write-Host "  git config http.postBuffer 524288000" -ForegroundColor White
Write-Host "  git push origin main" -ForegroundColor White
Write-Host ""
