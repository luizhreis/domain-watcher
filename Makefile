# DomainWatcher Makefile

.PHONY: build run test clean deps install dev help
.PHONY: test-unit test-integration test-storage test-memory test-sqlite test-postgresql
.PHONY: test-coverage test-coverage-html test-bench test-race
.PHONY: test-watch test-verbose test-short
.PHONY: postgresql-setup postgresql-test postgresql-cleanup

# Variáveis
BINARY_NAME=domain-watcher
BINARY_PATH=./bin/$(BINARY_NAME)
MAIN_PATH=./cmd/domain-watcher

# Cores para output
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
RED=\033[0;31m
NC=\033[0m # No Color

# Comando padrão
all: build

# Instalar dependências
deps:
	@echo "📦 Instalando dependências..."
	go mod download
	go mod tidy

# Build da aplicação
build: deps
	@echo "🔨 Compilando aplicação..."
	mkdir -p bin
	CGO_ENABLED=1 go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "✅ Build concluído: $(BINARY_PATH)"

# Executar aplicação
run: build
	@echo "🚀 Executando DomainWatcher..."
	$(BINARY_PATH)

# Executar em modo desenvolvimento (com rebuild automático)
dev:
	@echo "🔧 Modo desenvolvimento..."
	go run $(MAIN_PATH)

# ========================================
# TESTES
# ========================================

# TESTES - ORGANIZADOS POR CONTEXTO (WHITE-BOX vs BLACK-BOX)

# Executar todos os testes (unitários + integração)
test:
	@echo "$(BLUE)🧪 Executando todos os testes...$(NC)"
	@echo "$(GREEN)📦 Testes unitários (white-box):$(NC)"
	go test -short -v ./internal/...
	@echo "$(YELLOW)🔗 Testes de integração (black-box):$(NC)"
	go test -v ./tests/integration/...

# Testes unitários apenas (white-box - rápidos)
test-unit:
	@echo "$(GREEN)⚡ Testes unitários (white-box)...$(NC)"
	go test -short -v ./internal/...

# Testes de integração apenas (black-box - mais lentos)
test-integration:
	@echo "$(YELLOW)🔗 Testes de integração (black-box)...$(NC)"
	go test -v ./tests/integration/...

# Coverage completo (unitários + integração)
test-coverage:
	@echo "$(BLUE)📊 Coverage completo (unitários + integração)...$(NC)"
	@mkdir -p tests/coverage
	@echo "$(GREEN)📦 Coverage unitário (white-box):$(NC)"
	go test -coverprofile=tests/coverage/unit.out ./internal/...
	go tool cover -func=tests/coverage/unit.out
	@echo "$(YELLOW)🔗 Coverage integração (black-box):$(NC)"
	go test -coverprofile=tests/coverage/integration.out -coverpkg=./internal/... ./tests/integration/...
	go tool cover -func=tests/coverage/integration.out
	@echo "$(GREEN)✅ Coverage unitário: tests/coverage/unit.out$(NC)"
	@echo "$(GREEN)✅ Coverage integração: tests/coverage/integration.out$(NC)"
	@echo "$(BLUE)💡 Unitários medem implementação, integração mede comportamento$(NC)"

# Coverage apenas unitários (white-box - melhor coverage)
test-coverage-unit:
	@echo "$(BLUE)📊 Coverage unitários (white-box)...$(NC)"
	@mkdir -p tests/coverage
	go test -coverprofile=tests/coverage/unit.out ./internal/...
	go tool cover -func=tests/coverage/unit.out
	@echo "$(GREEN)✅ White-box coverage: tests/coverage/unit.out$(NC)"

# Coverage apenas integração (black-box)
test-coverage-integration:
	@echo "$(BLUE)📊 Coverage integração (black-box)...$(NC)"
	@mkdir -p tests/coverage
	go test -coverprofile=tests/coverage/integration.out -coverpkg=./internal/... ./tests/integration/...
	go tool cover -func=tests/coverage/integration.out
	@echo "$(YELLOW)✅ Black-box coverage: tests/coverage/integration.out$(NC)"

# Coverage detalhado com HTML
test-coverage-html:
	@echo "$(BLUE)📊 Coverage HTML (unitários)...$(NC)"
	@mkdir -p tests/coverage
	go test -coverprofile=tests/coverage/unit.out ./internal/...
	go tool cover -html=tests/coverage/unit.out -o tests/coverage/unit.html
	@echo "$(GREEN)✅ Coverage HTML: tests/coverage/unit.html$(NC)"
	@echo "$(BLUE)💡 Abra tests/coverage/unit.html no navegador$(NC)"

# Coverage HTML integração
test-coverage-html-integration:
	@echo "$(BLUE)📊 Coverage HTML integração...$(NC)"
	@mkdir -p tests/coverage
	go test -coverprofile=tests/coverage/integration.out -coverpkg=./internal/... ./tests/integration/...
	go tool cover -html=tests/coverage/integration.out -o tests/coverage/integration.html
	@echo "$(YELLOW)✅ Coverage HTML integração: tests/coverage/integration.html$(NC)"

# Coverage HTML combinado (experimental)
test-coverage-html-combined:
	@echo "$(BLUE)📊 Coverage HTML combinado (experimental)...$(NC)"
	@mkdir -p tests/coverage
	go test -coverprofile=tests/coverage/unit.out ./internal/...
	go test -coverprofile=tests/coverage/integration.out -coverpkg=./internal/... ./tests/integration/...
	@echo "mode: set" > tests/coverage/combined.out
	@grep -h -v "mode: set" tests/coverage/unit.out tests/coverage/integration.out >> tests/coverage/combined.out 2>/dev/null || true
	go tool cover -html=tests/coverage/combined.out -o tests/coverage/combined.html
	@echo "$(GREEN)✅ Coverage HTML combinado: tests/coverage/combined.html$(NC)"
	@echo "$(YELLOW)⚠️  Experimental - pode ter overlaps$(NC)"

# Todos os testes (unitários + integração)
test-all:
	@echo "$(BLUE)🚀 Todos os testes...$(NC)"
	@echo "$(YELLOW)⚠️  Incluindo SQLite (temporário)$(NC)"
	go test -v ./...

# Testes específicos de storage
test-storage:
	@echo "$(BLUE)🗄️ Testes de Storage...$(NC)"
	@echo "$(GREEN)Memory Storage (Production Ready):$(NC)"
	go test -v ./internal/storage/memory/tests
	@echo "$(YELLOW)SQLite Storage (Temporary - skipping race conditions):$(NC)"
	go test -v ./internal/storage/sqlite/tests -short

# Testes apenas Memory Storage
test-memory:
	@echo "$(GREEN)🧠 Memory Storage Tests...$(NC)"
	go test -v ./internal/storage/memory/tests

# Testes apenas SQLite Storage
test-sqlite:
	@echo "$(YELLOW)🗄️ SQLite Storage Tests...$(NC)"
	go test -v ./internal/storage/sqlite/tests

# Testes SQLite com integração
test-sqlite-full:
	@echo "$(YELLOW)🗄️ SQLite Storage Tests (com integração)...$(NC)"
	go test -v ./internal/storage/sqlite/tests

# Limpar arquivos de coverage
test-coverage-clean:
	@echo "$(YELLOW)🧹 Limpando coverage...$(NC)"
	rm -rf tests/coverage/
	@echo "$(GREEN)✅ Coverage limpo$(NC)"

# Coverage com limpeza automática
test-coverage-fresh:
	@echo "$(BLUE)📊 Coverage fresh (limpa + executa)...$(NC)"
	@make test-coverage-clean
	@make test-coverage
	@echo "$(GREEN)✅ Coverage fresh concluído!$(NC)"
