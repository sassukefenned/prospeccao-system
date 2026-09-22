.PHONY: help run build clean install deps reset-db

help:
	@echo "📋 Comandos disponíveis:"
	@echo "  make install    - Instalar dependências"
	@echo "  make run        - Executar servidor local (localhost:8080)"
	@echo "  make build      - Compilar para binário"
	@echo "  make clean      - Limpar arquivos de build"
	@echo "  make reset-db   - Resetar banco de dados (limpar tudo)"
	@echo "  make dev        - Executar em modo desenvolvimento"

install:
	@echo "📦 Instalando dependências..."
	@cd backend && go mod download

deps:
	@echo "📦 Baixando dependências..."
	@cd backend && go mod tidy

run:
	@echo "🚀 Iniciando servidor..."
	@cd backend && go run .

build:
	@echo "🔨 Compilando..."
	@cd backend && go build -o ../prospeccao .

clean:
	@echo "🧹 Limpando..."
	@rm -f prospeccao
	@rm -rf data/*.json data/*.csv
	@echo "✅ Limpeza concluída"

reset-db:
	@echo "⚠️  Deletando banco de dados..."
	@rm -f data/prospeccao.db
	@echo "✅ Banco de dados resetado"

dev:
	@echo "🛠️  Modo desenvolvimento..."
	@make reset-db
	@make run
