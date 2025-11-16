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

Projeto monolítico modular seguindo princípios de Clean Architecture:

- **Domain**: Entidades de negócio e regras de domínio
- **Use Cases**: Casos de uso e lógica de aplicação
- **Adapters**: Handlers HTTP e Repositories
- **Infrastructure**: Banco de dados, validadores, logger

### Módulos Disponíveis

1. **Auth**: Autenticação, autorização, usuários, roles e permissions
2. **Budget**: Contas, categorias, transações, orçamentos e relatórios

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

## 📚 Documentação da API

Após executar a aplicação, acesse:

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/health

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
