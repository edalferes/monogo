# Módulos

O Monetics é organizado em módulos independentes, cada um responsável por um domínio específico do negócio.

## Estrutura de um Módulo

Cada módulo segue a mesma estrutura padrão:

```
module/
├── domain/              # Entidades e regras de negócio
│   ├── entity.go
│   └── value_objects.go
├── usecase/            # Casos de uso
│   ├── interfaces/     # Interfaces de repositories
│   └── create_*.go
├── adapters/           # Adaptadores
│   ├── http/          # Handlers HTTP
│   │   ├── handlers/
│   │   ├── dto/
│   │   └── middleware.go
│   └── repository/    # Implementações de persistência
│       ├── gorm/
│       └── mappers/
├── service/           # Serviços auxiliares
├── errors/            # Erros específicos do módulo
├── entities.go        # Export de entidades para migrations
├── module.go          # WireUp do módulo
└── seed.go           # Dados iniciais
```

## Módulo Auth

Responsável por autenticação e autorização.

### Entidades

- **User**: Usuários do sistema
- **Role**: Papéis (admin, user, etc)
- **Permission**: Permissões granulares
- **AuditLog**: Log de auditoria

### Use Cases

- Criar/Listar/Atualizar/Deletar usuários
- Gerenciar roles e permissions
- Login e geração de JWT
- Verificar permissões
- Registrar auditoria

### Endpoints

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/v1/auth/login` | Login |
| POST | `/v1/auth/users` | Criar usuário |
| GET | `/v1/auth/users` | Listar usuários |
| GET | `/v1/auth/users/:id` | Buscar usuário |
| PUT | `/v1/auth/users/:id` | Atualizar usuário |
| DELETE | `/v1/auth/users/:id` | Deletar usuário |
| GET | `/v1/auth/roles` | Listar roles |
| POST | `/v1/auth/roles` | Criar role |
| GET | `/v1/auth/permissions` | Listar permissions |

### Serviços

- **JWTService**: Geração e validação de tokens
- **PasswordService**: Hash e verificação de senhas (bcrypt)
- **AuditService**: Registro de ações

### Middleware

- **AuthMiddleware**: Valida JWT
- **PermissionMiddleware**: Verifica permissões

## Módulo Budget

Responsável pela gestão financeira.

### Entidades

- **Account**: Contas bancárias
- **Category**: Categorias de transações
- **Transaction**: Transações financeiras
- **Budget**: Orçamentos planejados

### Use Cases

- CRUD de contas
- CRUD de categorias
- CRUD de transações
- CRUD de orçamentos
- Calcular saldo de conta
- Gerar relatórios mensais

### Endpoints

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/v1/budget/accounts` | Criar conta |
| GET | `/v1/budget/accounts` | Listar contas |
| GET | `/v1/budget/accounts/:id/balance` | Saldo da conta |
| POST | `/v1/budget/categories` | Criar categoria |
| GET | `/v1/budget/categories` | Listar categorias |
| POST | `/v1/budget/transactions` | Criar transação |
| GET | `/v1/budget/transactions` | Listar transações |
| POST | `/v1/budget/budgets` | Criar orçamento |
| GET | `/v1/budget/budgets` | Listar orçamentos |
| GET | `/v1/budget/reports/monthly` | Relatório mensal |

### Tipos de Transações

```go
const (
    TransactionTypeIncome   = "income"    // Receita
    TransactionTypeExpense  = "expense"   // Despesa
    TransactionTypeTransfer = "transfer"  // Transferência
)
```

### Tipos de Categorias

```go
const (
    CategoryTypeIncome  = "income"   // Receita
    CategoryTypeExpense = "expense"  // Despesa
)
```

### Períodos de Orçamento

```go
const (
    BudgetPeriodMonthly   = "monthly"    // Mensal
    BudgetPeriodQuarterly = "quarterly"  // Trimestral
    BudgetPeriodYearly    = "yearly"     // Anual
    BudgetPeriodCustom    = "custom"     // Personalizado
)
```

## Integração entre Módulos

### Auth → Budget

O módulo Budget depende do Auth para:

- Autenticação via JWT
- Identificação do usuário (UserID)
- Controle de acesso

Todas as rotas do Budget exigem autenticação.

### Isolamento

Cada módulo:

- Tem suas próprias entidades
- Define suas próprias interfaces
- Não depende diretamente de outros módulos
- Comunica-se via contratos bem definidos

## Adicionar um Novo Módulo

Para adicionar um novo módulo:

1. Crie a estrutura de diretórios
2. Defina as entidades de domínio
3. Implemente os use cases
4. Crie handlers e repositories
5. Implemente `module.go` com função `WireUp`
6. Registre em `internal/applications/api/app.go`

Exemplo de `module.go`:

```go
func WireUp(group *echo.Group, db *gorm.DB, jwtSecret string, logger logger.Logger) {
    // Inicializar repositories
    repo := repository.NewRepository(db)
    
    // Inicializar use cases
    useCase := usecase.NewUseCase(repo)
    
    // Inicializar handlers
    handler := handler.NewHandler(useCase)
    
    // Registrar rotas
    module := group.Group("/module")
    module.POST("", handler.Create)
    module.GET("", handler.List)
}
```

## Próximos Passos

- Veja o [Schema do Banco](database.md)
- Explore a [API do Auth](../api/auth.md)
- Explore a [API do Budget](../api/accounts.md)
