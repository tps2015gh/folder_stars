# Obsidian Star Graph - Interactive Launcher

Clear-Host
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "   Obsidian Star Graph Launcher" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan

$exe = ".\folder_star.exe"
if (-not (Test-Path $exe)) {
    Write-Host "Error: folder_star.exe not found. Please run build.bat first." -ForegroundColor Red
    exit
}

Write-Host "`nChoose a directory to scan:"
Write-Host " [0] Current Directory (.)" -ForegroundColor Yellow

$subdirs = Get-ChildItem -Directory | Where-Object { $_.Name -notlike ".*" }
$idx = 1
foreach ($dir in $subdirs) {
    Write-Host " [$idx] $($dir.Name)"
    $idx++
}
Write-Host " [C] Enter Custom Path"

$choice = Read-Host "`nSelection (default 0)"
if ([string]::IsNullOrWhiteSpace($choice)) { $choice = "0" }

$selectedDir = "."
if ($choice -eq "C") {
    $selectedDir = Read-Host "Enter custom path"
} elseif ($choice -match "^\d+$") {
    $num = [int]$choice
    if ($num -gt 0 -and $num -lt $idx) {
        $selectedDir = $subdirs[$num-1].Name
    }
}

$port = Read-Host "Enter port (default 8080)"
if ([string]::IsNullOrWhiteSpace($port)) { $port = "8080" }

Write-Host "`nStarting server for: $selectedDir on port $port..." -ForegroundColor Green
& $exe -dir $selectedDir -port $port
