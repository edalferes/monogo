# Instalação

Este guia mostra como instalar e executar o Monetics em diferentes ambientes.

## Pré-requisitos

Antes de começar, certifique-se de ter instalado:

- **Go 1.25.1+** - [Download](https://golang.org/dl/)
- **Docker e Docker Compose** - [Download](https://www.docker.com/get-started)
- **Make** - Geralmente já instalado em sistemas Unix
- **PostgreSQL 14+** - Necessário se não usar Docker

## Instalação com Docker Compose

A forma mais rápida de executar o Monetics é usando Docker Compose:

```bash
# Clone o repositório
git clone https://github.com/alpheres/monetics.git
cd monetics

# Suba os containers
docker-compose up -d

# Verifique se está funcionando
curl http://localhost:8080/health
```

A API estará disponível em `http://localhost:8080`.

## Instalação Local

### 1. Clone o Repositório

```bash
git clone https://github.com/alpheres/monetics.git
cd monetics
```

### 2. Instale as Dependências

```bash
go mod download
```

### 3. Configure o Banco de Dados

Certifique-se de ter um PostgreSQL rodando e configure as credenciais no arquivo `config.yaml`:

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  name: monetics
  sslmode: disable
```

### 4. Execute a Aplicação

```bash
# Usando Make
make run

# Ou compilando primeiro
make build
./bin/monetics
```

## Verificação da Instalação

Após iniciar a aplicação, você pode verificar se tudo está funcionando:

```bash
# Health check
curl http://localhost:8080/health

# Swagger UI
open http://localhost:8080/swagger/index.html
```

## Credenciais Padrão

Um usuário root é criado automaticamente no primeiro start:

- **Username**: `root`
- **Password**: `root123`

!!! warning "Produção"
    Lembre-se de alterar as credenciais padrão em ambientes de produção!

## Próximos Passos

- Configure as variáveis de ambiente em [Configuração](configuration.md)
- Aprenda a usar a API em [Primeiro Uso](quickstart.md)
- Explore os endpoints na [API Reference](../api/auth.md)
