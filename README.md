# Monetics

**Monetics** é uma API RESTful para gestão financeira pessoal, desenvolvida em Go com arquitetura modular e foco em boas práticas de desenvolvimento.

## 📋 Sobre o Projeto

O Monetics oferece um sistema completo para controle de finanças pessoais, permitindo:

- **Gestão de Contas**: Criação e gerenciamento de contas bancárias com saldo atualizado
- **Categorização**: Organização de receitas e despesas em categorias personalizadas
- **Transações**: Registro de receitas, despesas e transferências entre contas
- **Orçamentos**: Planejamento financeiro com definição de limites por categoria e período
- **Relatórios**: Visualização de gastos mensais e acompanhamento de orçamentos
- **Autenticação e Autorização**: Sistema completo com roles e permissions baseado em JWT
- **Auditoria**: Registro de ações dos usuários para rastreabilidade

## 🚀 Tecnologias Utilizadas

- **Go 1.25.1** - Linguagem de programação
- **Echo Framework v4** - Framework web para rotas HTTP
- **GORM** - ORM para persistência de dados
- **PostgreSQL** - Banco de dados relacional
- **JWT** - Autenticação via tokens
- **Swagger** - Documentação automática da API
- **Docker & Docker Compose** - Containerização
- **Zerolog** - Logging estruturado

## 🏗️ Arquitetura

Projeto modular inspirado na arquitetura do **Grafana Loki**, oferecendo flexibilidade para executar como monólito ou microserviços.

### Princípios

- **Clean Architecture**: Separação clara de Domain, Use Cases, Adapters e Infrastructure
- **Modularidade**: Cada módulo pode ser executado independentemente
- **Comunicação HTTP**: Inspirada no Loki, com suporte a comunicação local e remota
- **Dependency Injection**: Container e Registry para gerenciar dependências entre módulos
- **Health Checks**: Endpoints `/health`, `/ready` e `/live` para Kubernetes

### Módulos Disponíveis

1. **Auth**: Autenticação, autorização, usuários, roles e permissions
2. **Budget**: Contas, categorias, transações, orçamentos e relatórios (depende de Auth)

### Modos de Execução

```bash
# Executar todos os módulos (monólito)
./bin/monetics --module=all

# Executar apenas auth (microservice)
./bin/monetics --module=auth

# Executar apenas budget (microservice)
./bin/monetics --module=budget

# Executar múltiplos módulos
./bin/monetics --module=auth,budget
```

### Sistema de Dependências

O Monetics usa um **Dependency Injection Container** e **Module Registry** para gerenciar dependências entre módulos:

- **Inicialização automática**: Módulos são inicializados na ordem correta baseado no grafo de dependências
- **Comunicação flexível**: Usa serviços locais (in-memory) quando disponíveis, caso contrário HTTP
- **Fail-fast**: Valida dependências no startup, evitando erros em runtime

**Exemplo de dependência**:
```go
// Budget depende de Auth
registry.Register("budget", []string{"auth"}, ...)
```

**Monólito** (`--module=all`):
- Auth inicializa primeiro
- Budget usa serviço local do Auth (in-memory)
- Zero overhead de rede

**Microservices** (módulos separados):
- Cada módulo roda em processo independente
- Budget conecta via HTTP ao Auth service
- Retry + Circuit Breaker automático

📚 **Documentação completa**: 
- [Module Dependencies Guide](./docs/module-dependencies.md)
- [Testing Dependencies](./docs/testing-dependencies.md)

### Comunicação Entre Serviços

A comunicação entre módulos é configurada via `config.yaml`:

```yaml
modules:
  auth:
    url: ""  # Vazio = local (monólito)
    # url: "http://auth-service:8080"  # HTTP remoto (microservices)
```

O sistema usa **pkg/httpclient** para comunicação HTTP resiliente com:

- Connection pooling otimizado (100 conexões idle)
- Timeouts configuráveis (10s padrão)
- Retry logic com backoff exponencial (3 tentativas)
- Circuit breaker para proteção contra falhas em cascata
- Jitter aleatório para evitar thundering herd

## 📦 Como Executar

### Pré-requisitos

- Go 1.25.1+
- Docker e Docker Compose
- Make

### Executando com Docker Compose

```sh
# Subir banco de dados e aplicação
docker-compose up

# A API estará disponível em http://localhost:8080
```

### Executando localmente

```sh
# Instalar dependências
go mod download

# Executar a aplicação
make run

# Ou compilar primeiro
make build
./bin/monetics
```

## 📚 Documentação

O projeto possui documentação completa organizada por tipo:

### 📘 Developer Documentation (MkDocs)

Documentação para desenvolvedores com guias, arquitetura e tutoriais:

```bash
# Servir localmente em http://127.0.0.1:8000
mkdocs serve
```

**Conteúdo**:
- [Module Dependencies](./docs/mkdocs/module-dependencies.md) - Sistema de injeção de dependências
- [Testing Guide](./docs/mkdocs/testing-dependencies.md) - Como testar dependências
- [Dependency Graph](./docs/mkdocs/dependency-graph.md) - Visualização de dependências
- [Architecture](./docs/mkdocs/architecture/) - Documentação de arquitetura
- [Getting Started](./docs/mkdocs/getting-started/) - Guias de início

### 🔧 API Documentation (Swagger)

Documentação interativa da API REST:

```bash
# Gerar/atualizar Swagger docs
make swagger
```

**Acesso**: http://localhost:8080/swagger/index.html (quando API estiver rodando)

### 🧪 Postman Collection

Collection completa para testes manuais:

**Arquivo**: [`docs/postman/Monetics.postman_collection.json`](./docs/postman/Monetics.postman_collection.json)

**Como usar**:
1. Importar no Postman
2. Configurar variáveis de ambiente (BASE_URL, AUTH_URL, BUDGET_URL)
3. Testar endpoints

📖 **Mais detalhes**: [docs/README.md](./docs/README.md)

### Credenciais Padrão

Usuário root criado automaticamente no seed:
- **Username**: `root`
- **Password**: `root123`

## 🧪 Testes

```sh
make test
```

## 📁 Estrutura de Pastas

```
├── cmd/                    # Entrypoint da aplicação
│   └── api/
├── internal/               # Código principal da aplicação
│   ├── applications/       # Configuração da aplicação
│   ├── config/            # Configurações e variáveis de ambiente
│   ├── infra/             # Infraestrutura (DB, validators)
│   └── modules/           # Módulos de domínio
│       ├── auth/          # Autenticação e autorização
│       └── budget/        # Gestão financeira
├── pkg/                   # Pacotes reutilizáveis
│   ├── logger/            # Logger configurável
│   └── responses/         # Respostas HTTP padronizadas
├── docs/                  # Documentação Swagger
└── scripts/               # Scripts auxiliares
```

## 🔧 Comandos Make Disponíveis

```sh
make help         # Exibir comandos disponíveis
make build        # Compilar a aplicação
make run          # Executar a aplicação
make test         # Executar testes
make swagger      # Gerar documentação Swagger
make clean        # Limpar artifacts de build
make docker-build # Construir imagem Docker
```

## 🌟 Principais Características

- ✅ Arquitetura modular e escalável
- ✅ Separação clara de responsabilidades (Clean Architecture)
- ✅ Documentação automática com Swagger
- ✅ Validação de dados com go-playground/validator
- ✅ Respostas HTTP padronizadas
- ✅ Sistema completo de autenticação e autorização
- ✅ Logging estruturado com níveis configuráveis
- ✅ Migrations automáticas com GORM
- ✅ Seed de dados iniciais
- ✅ Pronto para containerização
- ✅ Integração com Backstage TechDocs

## 📖 Documentação Técnica

A documentação completa do projeto está disponível via **Backstage TechDocs** usando MkDocs.

### Visualizar Localmente

```bash
# Instalar MkDocs
pip install mkdocs-techdocs-core

# Servir documentação
mkdocs serve

# Acesse: http://localhost:8000
```

### Backstage Integration

O projeto está configurado para integração com o Backstage:

- **Catalog**: `catalog-info.yaml` - Define componentes e APIs
- **TechDocs**: `mkdocs.yml` + `docs/` - Documentação técnica completa

Consulte [docs/README.md](docs/README.md) para mais detalhes sobre a documentação.
