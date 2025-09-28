# Configuração

Este projeto suporta múltiplas formas de configuração para flexibilidade em diferentes ambientes.

## Ordem de Precedência

1. **Variáveis de ambiente** (maior prioridade)
2. **Arquivos YAML** (config.yaml, config.dev.yaml, etc.)
3. **Valores padrão** (menor prioridade)

## Arquivos de Configuração

### config.yaml
Configuração base com valores padrão para desenvolvimento.

### config.dev.yaml
Configuração específica para desenvolvimento local.

### config.prod.yaml
Configuração para produção (valores sensíveis devem vir de env vars).

## Variáveis de Ambiente

Todas as configurações podem ser sobrescritas via variáveis de ambiente usando o padrão:
- `APP_NAME` → `app.name`
- `DATABASE_HOST` → `database.host`
- `JWT_SECRET` → `jwt.secret`

Veja `.env.example` para uma lista completa.

## Uso em Código

```go
import "github.com/edalferes/monogo/pkg/config"

// Carrega configuração automaticamente
cfg := config.LoadConfig()

// Ou com opções específicas
cfg, err := config.Load(config.ConfigOptions{
    ConfigPath: "./configs",
    ConfigName: "app",
    ConfigType: "yaml",
})

// Métodos utilitários
fmt.Println(cfg.GetDSN())        // String de conexão do DB
fmt.Println(cfg.IsDevelopment()) // true se environment == "development"
fmt.Println(cfg.IsProduction())  // true se environment == "production"
```

## Validação

A configuração é automaticamente validada no carregamento:
- Campos obrigatórios
- Valores de porta válidos
- Ambiente válido (development, staging, production)

## Exemplo de Uso

```bash
# Usando arquivo YAML
./app

# Usando variáveis de ambiente
DATABASE_HOST=prod-db.example.com \
JWT_SECRET=super-secret-key \
./app

# Carregando arquivo específico
CONFIG_PATH=./configs CONFIG_NAME=production ./app
```