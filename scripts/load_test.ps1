# scripts/load_test.ps1 — gRPC load test using ghz (Windows PowerShell)
#
# Usage:
#   .\scripts\load_test.ps1
#   .\scripts\load_test.ps1 -Total 50000 -Concurrency 200 -Addr localhost:50051

param(
    [string]$Addr        = "localhost:50051",
    [int]   $Total       = 10000,
    [int]   $Concurrency = 100,
    [int]   $Timeout     = 30
)

Write-Host "=== Distributed Job Queue — Load Test ===" -ForegroundColor Cyan
Write-Host "  Server      : $Addr"
Write-Host "  Total jobs  : $Total"
Write-Host "  Concurrency : $Concurrency"
Write-Host "  Timeout     : ${Timeout}s"
Write-Host ""

# Check ghz is installed
if (-not (Get-Command ghz -ErrorAction SilentlyContinue)) {
    Write-Host "[!] ghz not found. Install it with:" -ForegroundColor Red
    Write-Host "    go install github.com/bojand/ghz/cmd/ghz@latest" -ForegroundColor Yellow
    exit 1
}

# Write payload to a temp file (avoids Windows quote-stripping on native executables)
$payload = '{"type":"echo","payload":"aGVsbG8gd29ybGQ=","priority":0,"delay_seconds":0,"max_retries":3}'
$tmpFile = [System.IO.Path]::GetTempFileName()
Set-Content -Path $tmpFile -Value $payload -Encoding utf8

Write-Host "Starting load test..." -ForegroundColor Green
Write-Host ""

ghz `
    --insecure `
    --proto "api/proto/job.proto" `
    --call "jobpb.JobService.SubmitJob" `
    --data-file $tmpFile `
    --total $Total `
    --concurrency $Concurrency `
    --timeout "${Timeout}s" `
    $Addr

Remove-Item $tmpFile -Force
