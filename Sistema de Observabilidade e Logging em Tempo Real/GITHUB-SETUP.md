# GitHub Setup Guide

## Problem: Files too large for GitHub

GitHub has limits:
- **100MB** per file
- **1GB** recommended repository size
- **5GB** hard limit

## Solution

### Step 1: Clean Repository

```powershell
.\scripts\prepare-for-github.ps1
```

This removes:
- Build artifacts (bin/, *.exe)
- Logs (*.log, logs/)
- Cache files (tmp/, .cache/)
- Temporary files

### Step 2: Increase Git Buffer

```bash
git config http.postBuffer 524288000
git config https.postBuffer 524288000
```

### Step 3: Commit in Chunks

```powershell
.\scripts\commit-in-chunks.ps1
```

This commits files in smaller groups:
1. Core Go files
2. Web files
3. Scripts
4. Configuration
5. Documentation
6. Infrastructure
7. Proto files

### Step 4: Push to GitHub

```bash
git push origin main
```

## Alternative: Use Git LFS

For large binary files:

```bash
git lfs install
git lfs track "*.exe"
git lfs track "*.dll"
git add .gitattributes
git commit -m "Add Git LFS"
```

## Troubleshooting

### Error: "file too large"

1. Check file sizes:
   ```powershell
   .\scripts\check-file-sizes.ps1
   ```

2. Add large files to `.gitignore`

3. Remove from Git history:
   ```bash
   git rm --cached path/to/large/file
   ```

### Error: "pack exceeds maximum allowed size"

Use chunked commits:
```powershell
.\scripts\commit-in-chunks.ps1
```

### Error: "remote: error: GH001"

Repository too large. Options:
1. Use Git LFS
2. Split into multiple repositories
3. Remove unnecessary files

## Best Practices

1. **Never commit:**
   - Build artifacts (bin/, *.exe)
   - Dependencies (node_modules/, vendor/)
   - Logs (*.log)
   - Environment files (.env)
   - Large binaries

2. **Always use .gitignore**

3. **Keep repository under 1GB**

4. **Use Git LFS for:**
   - Videos
   - Large images
   - Binary assets
   - Compiled files

## Quick Commands

```bash
# Check repository size
git count-objects -vH

# Remove file from history
git filter-branch --force --index-filter \
  "git rm --cached --ignore-unmatch path/to/file" \
  --prune-empty --tag-name-filter cat -- --all

# Clean up
git reflog expire --expire=now --all
git gc --prune=now --aggressive
```
