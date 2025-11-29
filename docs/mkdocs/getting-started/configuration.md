# Configuração

O Monetics utiliza o [Viper](https://github.com/spf13/viper) para gerenciamento de configurações, suportando múltiplos formatos e variáveis de ambiente.

## Arquivos de Configuração

O projeto vem com três arquivos de configuração:

- `config.yaml` - Configuração base
- `config.dev.yaml` - Configurações de desenvolvimento
- `config.prod.yaml` - Configurações de produção

## Estrutura de Configuração

### Aplicação

```yaml
app:
  name: monetics
  port: 8080
  env: development
```

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `name` | string | Nome da aplicação |
| `port` | int | Porta HTTP |
| `env` | string | Ambiente (development, production) |

### Banco de Dados

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  name: monetics
  sslmode: disable
  max_connections: 100
  max_idle_connections: 10
  connection_max_lifetime: 3600
```

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `host` | string | Host do PostgreSQL |
| `port` | int | Porta do PostgreSQL |
| `user` | string | Usuário do banco |
| `password` | string | Senha do banco |
| `name` | string | Nome do banco |
| `sslmode` | string | Modo SSL (disable, require) |
| `max_connections` | int | Máximo de conexões |
| `max_idle_connections` | int | Conexões idle |
| `connection_max_lifetime` | int | Tempo de vida (segundos) |

### JWT

```yaml
jwt:
  secret: your-secret-key-change-in-production
  expiration: 24h
```

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `secret` | string | Chave secreta para assinar tokens |
| `expiration` | duration | Tempo de expiração do token |

!!! danger "Segurança"
    **NUNCA** commite o `jwt.secret` em produção! Use variáveis de ambiente.

### Logger

```yaml
logger:
  level: info
  format: json
```

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `level` | string | Nível de log (debug, info, warn, error) |
| `format` | string | Formato (json, console) |

## Variáveis de Ambiente

Todas as configurações podem ser sobrescritas por variáveis de ambiente usando o prefixo `MONETICS_`:

```bash
export MONETICS_APP_PORT=9090
export MONETICS_DATABASE_HOST=db.example.com
export MONETICS_DATABASE_PASSWORD=secret
export MONETICS_JWT_SECRET=super-secret-key
export MONETICS_LOGGER_LEVEL=debug
```

## Configuração Docker

No `docker-compose.yml`, as variáveis são definidas:

```yaml
environment:
  - MONETICS_DATABASE_HOST=postgres
  - MONETICS_DATABASE_PORT=5432
  - MONETICS_DATABASE_USER=postgres
  - MONETICS_DATABASE_PASSWORD=postgres
  - MONETICS_DATABASE_NAME=monetics
  - MONETICS_JWT_SECRET=change-this-secret-in-production
```

## Boas Práticas

1. **Desenvolvimento**: Use `config.dev.yaml` com valores de teste
2. **Produção**: Use variáveis de ambiente para valores sensíveis
3. **Secrets**: Nunca commite senhas ou chaves em arquivos de configuração
4. **Validação**: O Monetics valida todas as configurações no startup

## Exemplo Completo

```yaml
app:
  name: monetics
  port: 8080
  env: production

database:
  host: db.prod.example.com
  port: 5432
  user: monetics_user
  password: ${DB_PASSWORD}  # Use variável de ambiente
  name: monetics_prod
  sslmode: require
  max_connections: 50
  max_idle_connections: 5
  connection_max_lifetime: 1800

jwt:
  secret: ${JWT_SECRET}  # Use variável de ambiente
  expiration: 8h

logger:
  level: info
  format: json
```

## Próximos Passos

- Aprenda a usar a API em [Primeiro Uso](quickstart.md)
- Entenda a arquitetura em [Visão Geral](../architecture/overview.md)
