# Visão Geral da Arquitetura

O Monetics segue os princípios de **Clean Architecture** e **Domain-Driven Design (DDD)**, garantindo separação de responsabilidades e facilidade de manutenção.

## Princípios Arquiteturais

### 1. Independência de Frameworks
O core da aplicação não depende de frameworks específicos. Echo, GORM e outras bibliotecas são apenas detalhes de implementação que podem ser substituídos.

### 2. Testabilidade
A separação em camadas permite testar a lógica de negócio sem depender de infraestrutura (banco de dados, HTTP, etc).

### 3. Independência de UI
A API é independente de interfaces. Poderia ser consumida por web, mobile, CLI, etc.

### 4. Independência de Banco de Dados
O domínio não conhece o banco de dados. Repositories abstraem a persistência.

## Camadas da Aplicação

```
┌─────────────────────────────────────────┐
│          HTTP Handlers (Echo)           │
│         (Adapters/Interface)            │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│            Use Cases                    │
│       (Application Business Rules)      │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│            Domain Entities              │
│       (Enterprise Business Rules)       │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│          Repositories (GORM)            │
│         (Adapters/Database)             │
└─────────────────────────────────────────┘
```

### Domain (Entities)

Camada mais interna, contém:

- Entidades de negócio puras
- Regras de negócio empresariais
- Independente de qualquer framework

**Exemplo**: `User`, `Transaction`, `Budget`

```go
type Transaction struct {
    ID          uint
    UserID      uint
    AccountID   uint
    CategoryID  uint
    Type        TransactionType
    Amount      float64
    Description string
    Date        time.Time
    Status      TransactionStatus
}
```

### Use Cases

Camada de aplicação, contém:

- Casos de uso específicos do sistema
- Orquestração de entidades
- Regras de negócio da aplicação

**Exemplo**: `CreateTransaction`, `GetMonthlyReport`

```go
type CreateTransactionUseCase struct {
    transactionRepo TransactionRepository
    accountRepo     AccountRepository
}

func (uc *CreateTransactionUseCase) Execute(input CreateTransactionInput) error {
    // Validações
    // Regras de negócio
    // Persistência via repository
}
```

### Adapters

Camada de adaptadores, contém:

- **HTTP Handlers**: Controladores REST (Echo)
- **Repositories**: Implementações de persistência (GORM)
- **DTOs**: Data Transfer Objects

**Exemplo**: Handler HTTP

```go
func (h *TransactionHandler) Create(c echo.Context) error {
    var dto CreateTransactionDTO
    if err := c.Bind(&dto); err != nil {
        return err
    }
    
    result, err := h.useCase.Execute(dto)
    return c.JSON(http.StatusCreated, result)
}
```

### Infrastructure

Camada mais externa, contém:

- Configuração de banco de dados
- Logger
- Validadores
- Frameworks

## Estrutura de Diretórios

```
monetics/
├── cmd/                    # Entrypoints
│   └── api/
│       └── main.go
├── internal/              # Código da aplicação
│   ├── applications/      # Configuração da app
│   ├── config/           # Gerenciamento de configs
│   ├── infra/            # Infraestrutura
│   │   ├── db/
│   │   └── validator/
│   └── modules/          # Módulos de domínio
│       ├── auth/         # Autenticação
│       │   ├── domain/          # Entidades
│       │   ├── usecase/         # Casos de uso
│       │   └── adapters/        # HTTP & Repos
│       └── budget/       # Gestão financeira
│           ├── domain/
│           ├── usecase/
│           └── adapters/
└── pkg/                  # Pacotes compartilhados
    ├── logger/
    └── responses/
```

## Fluxo de Dados

1. **Request HTTP** chega no Handler (Echo)
2. **Handler** valida e converte para DTO
3. **Use Case** recebe DTO e executa lógica
4. **Domain** aplica regras de negócio
5. **Repository** persiste dados (GORM)
6. **Response** retorna pelo Handler

## Vantagens desta Arquitetura

✅ **Testabilidade**: Cada camada pode ser testada isoladamente  
✅ **Manutenibilidade**: Mudanças em uma camada não afetam outras  
✅ **Flexibilidade**: Fácil trocar frameworks ou banco de dados  
✅ **Escalabilidade**: Módulos podem virar microsserviços  
✅ **Clareza**: Código organizado e fácil de entender  

## Próximos Passos

- Entenda os [Módulos](modules.md) do sistema
- Veja o [Schema do Banco de Dados](database.md)
- Explore a [API Reference](../api/auth.md)
