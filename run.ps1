# Script PowerShell para rodar o Prospeccao System
# Uso: .\run.ps1

Write-Host "🚀 Iniciando Prospeccao System..." -ForegroundColor Green

# Verifica se Go está instalado
$goCheck = go version 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Go não está instalado!" -ForegroundColor Red
    Write-Host "Baixe em: https://golang.org/dl" -ForegroundColor Yellow
    exit 1
}

Write-Host "✅ Go detectado: $goCheck" -ForegroundColor Green

# Download das dependências
Write-Host "`n📦 Baixando dependências..." -ForegroundColor Cyan
go mod download
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Erro ao baixar dependências" -ForegroundColor Red
    exit 1
}

Write-Host "✅ Dependências prontas" -ForegroundColor Green

# Inicia o servidor
Write-Host "`n🔥 Iniciando servidor na porta 8080..." -ForegroundColor Cyan
Write-Host "📍 Acesse: http://localhost:8080" -ForegroundColor Yellow
Write-Host "Para parar: Ctrl + C`n" -ForegroundColor Gray

# Aguarda 2 segundos e abre o navegador
Start-Sleep -Seconds 2
Start-Process "http://localhost:8080"

# Roda o servidor
go run ./backend
