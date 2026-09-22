# 🎯 Sistema de Prospectação de Leads

Sistema **100% ISOLADO** para prospectação de leads usando Google Maps Scraper. Roda localmente com interface web simples.

## ✨ Características

- ✅ **Standalone** - Funciona completamente independente
- ✅ **Interface Web** - Dashboard intuitivo
- ✅ **SQLite** - Banco de dados local (sem dependências)
- ✅ **REST API** - Para integração futura
- ✅ **Export** - CSV + JSON com contatos inclusos
- ✅ **Campanhas** - Gerenciar prospecção por tipo
- ✅ **Filtros Avançados** - Por qualidade, status, perfil
- ✅ **Contatos Automáticos** - Email + Telefone em exportações

---

## 🚀 Começar Rápido

### 1️⃣ Pré-requisitos

- **Go 1.21+** ([Download](https://golang.org/dl/))
- **Terminal** (Bash, PowerShell, Zsh)

### 2️⃣ Setup

```bash
# Clonar/extrair o projeto
cd prospeccao-system

# Criar arquivo .env
cp .env.example .env

# Instalar dependências
make install
```

### 3️⃣ Executar

```bash
make run
```

Acesse: **http://localhost:8080**

---

## 📖 Como Usar

### Dashboard
- Ver estatísticas gerais de leads
- Taxa de conversão
- Distribuição por perfil e categoria

### Nova Busca
1. Acesse a aba "Nova Busca"
2. Digite o termo (ex: "distribuidora Goiânia")
3. Selecione categoria (opcional)
4. Clique "Iniciar Busca"

*Nota: Atualmente, a integração com Google Maps Scraper está em desenvolvimento*

### Gerenciar Leads
- Filtrar por status, qualidade, perfil
- Mudar status de leads (novo → contato → conversão)
- Deletar leads indesejados
- Adicionar notas

### Campanhas
- Criar campanhas por tipo de cliente
- Acompanhar conversões
- Medir taxa de sucesso

### Exportar
- **CSV**: Compatível com Excel
- **JSON**: Inclui contatos (email + telefone)

---

## 📋 Estrutura

```
prospeccao-system/
├── backend/
│   ├── main.go          # Entrada
│   ├── models.go        # Estruturas
│   ├── database.go      # SQLite
│   └── handlers.go      # API endpoints
├── frontend/
│   ├── index.html       # Interface
│   └── assets/
│       ├── app.js       # Lógica
│       └── style.css    # Estilos
├── data/                # Banco de dados + exports
├── go.mod               # Dependências
├── Makefile             # Comandos úteis
└── .env.example         # Config
```

---

## 🔧 Comandos

```bash
# Desenvolvimento
make run            # Iniciar servidor
make build          # Compilar binário
make clean          # Limpar arquivos
make reset-db       # Limpar BD completamente
make dev            # Modo dev (limpa BD + inicia)

# Instalação
make install        # Instalar dependências
make deps           # Atualizar dependências
```

---

## 🗄️ Banco de Dados

**SQLite Local** em `./data/prospeccao.db`

Tabelas:
- `leads` - Leads qualificados
- `buscas` - Histórico de buscas
- `campanhas` - Campanhas de prospecção

Criadas automaticamente na primeira execução.

---

## 🔌 Integração com Google Maps Scraper

A integração está estruturada em `handlers.go`:

```go
func executeSearch(search *Search) {
    // Aqui integra com: ./gmapsscraper search -q "termo"
    // Processa JSON de resultado
    // Qualifica leads (telefone, email, website)
    // Salva em BD
}
```

Próximos passos:
1. Colocar `gmapsscraper` no PATH
2. Implementar parsing do JSON de resultado
3. Aplicar filtros de qualificação
4. Atualizar status da busca

---

## 📧 Contatos Exportados

Todos os exports incluem:

```json
{
    "timestamp": "2024-09-22T21:00:00Z",
    "contato": {
        "email": "stfenned@gmail.com",
        "phone": "+55 62 8178-9507"
    },
    "leads": [...]
}
```

---

## 🎯 Próximas Versões

- [ ] Integração real com Google Maps Scraper
- [ ] Agendamento automático de buscas
- [ ] Machine learning para qualificação
- [ ] Integração com WhatsApp/SMS
- [ ] Docker para deploy em VPS

---

## ⚙️ Troubleshooting

### Porta 8080 já está em uso
```bash
# Mudar em .env
PORT=8081
```

### Erro ao compilar
```bash
cd backend
go mod tidy
go mod download
```

### Banco de dados corrompido
```bash
make reset-db
make run
```

---

## 📝 Licença

Privado - Projeto Victoria

---

## 💡 Suporte

Para dúvidas ou problemas:
- 📧 stfenned@gmail.com
- 📱 +55 62 8178-9507

---

**🚀 Pronto para prospeccionar!**
