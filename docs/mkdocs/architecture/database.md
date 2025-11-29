# Banco de Dados

O Monetics utiliza **PostgreSQL** como banco de dados relacional, com GORM para ORM e migrations automáticas.

## Visão Geral do Schema

```
┌─────────────┐       ┌──────────────┐       ┌─────────────┐
│    users    │──────<│  user_roles  │>──────│    roles    │
└─────────────┘       └──────────────┘       └─────────────┘
                                                     │
                                                     │
                                              ┌──────▼───────┐
                                              │ role_perms   │
                                              └──────┬───────┘
                                                     │
                                              ┌──────▼───────┐
                                              │ permissions  │
                                              └──────────────┘

┌─────────────┐       ┌──────────────┐       ┌─────────────┐
│   accounts  │<──────│ transactions │──────>│ categories  │
└─────────────┘       └──────────────┘       └─────────────┘
      │                      │                       │
      │                      └───────────────────────┘
      │                                              │
      └──────────────────────────────────────────────┘
                             │
                      ┌──────▼───────┐
                      │   budgets    │
                      └──────────────┘
```

## Tabelas do Módulo Auth

### users

Armazena os usuários do sistema.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | SERIAL | PK |
| username | VARCHAR(50) | UNIQUE, NOT NULL |
| password_hash | VARCHAR(255) | NOT NULL |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

### roles

Define os papéis no sistema.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | SERIAL | PK |
| name | VARCHAR(50) | UNIQUE, NOT NULL |
| description | TEXT | Descrição do role |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

**Roles padrão**: `admin`, `user`

### permissions

Define permissões granulares.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | SERIAL | PK |
| name | VARCHAR(100) | UNIQUE, NOT NULL |
| description | TEXT | Descrição |
| resource | VARCHAR(50) | Recurso (users, accounts, etc) |
| action | VARCHAR(50) | Ação (create, read, update, delete) |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

### user_roles

Relacionamento muitos-para-muitos entre users e roles.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| user_id | INTEGER | FK → users.id |
| role_id | INTEGER | FK → roles.id |

**PK Composta**: (user_id, role_id)

### role_permissions

Relacionamento muitos-para-muitos entre roles e permissions.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| role_id | INTEGER | FK → roles.id |
| permission_id | INTEGER | FK → permissions.id |

**PK Composta**: (role_id, permission_id)

### audit_logs

Registra ações dos usuários para auditoria.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | SERIAL | PK |
| user_id | INTEGER | FK → users.id, NOT NULL |
| action | VARCHAR(50) | Ação realizada |
| resource_type | VARCHAR(50) | Tipo do recurso |
| resource_id | INTEGER | ID do recurso |
| ip_address | VARCHAR(45) | IP do usuário |
| details | TEXT | Detalhes adicionais |
| created_at | TIMESTAMP | NOT NULL |

## Tabelas do Módulo Budget

### accounts

Contas bancárias dos usuários.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | SERIAL | PK |
| user_id | INTEGER | FK → users.id, NOT NULL |
| name | VARCHAR(100) | NOT NULL |
| type | VARCHAR(20) | checking, savings, investment |
| balance | DECIMAL(15,2) | Saldo atual |
| currency | VARCHAR(3) | BRL, USD, EUR |
| is_active | BOOLEAN | Conta ativa? |
| description | TEXT | Descrição |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

**Índices**: user_id, is_active

### categories

Categorias para organizar transações.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | SERIAL | PK |
| user_id | INTEGER | FK → users.id, NOT NULL |
| name | VARCHAR(100) | NOT NULL |
| type | VARCHAR(20) | income, expense |
| color | VARCHAR(7) | Cor (hex) |
| icon | VARCHAR(50) | Nome do ícone |
| parent_id | INTEGER | FK → categories.id (subcategorias) |
| is_active | BOOLEAN | Categoria ativa? |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

**Índices**: user_id, type, parent_id

### transactions

Transações financeiras (receitas, despesas, transferências).

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | SERIAL | PK |
| user_id | INTEGER | FK → users.id, NOT NULL |
| account_id | INTEGER | FK → accounts.id, NOT NULL |
| category_id | INTEGER | FK → categories.id |
| type | VARCHAR(20) | income, expense, transfer |
| amount | DECIMAL(15,2) | Valor, NOT NULL |
| description | TEXT | Descrição |
| date | TIMESTAMP | Data da transação, NOT NULL |
| status | VARCHAR(20) | pending, completed, cancelled |
| to_account_id | INTEGER | FK → accounts.id (para transfers) |
| tags | VARCHAR(255) | Tags separadas por vírgula |
| notes | TEXT | Notas adicionais |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

**Índices**: user_id, account_id, category_id, date, type, status

### budgets

Orçamentos planejados por categoria.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | SERIAL | PK |
| user_id | INTEGER | FK → users.id, NOT NULL |
| category_id | INTEGER | FK → categories.id, NOT NULL |
| name | VARCHAR(100) | NOT NULL |
| amount | DECIMAL(15,2) | Valor planejado, NOT NULL |
| spent | DECIMAL(15,2) | Valor gasto (calculado) |
| period | VARCHAR(20) | monthly, quarterly, yearly, custom |
| start_date | TIMESTAMP | Data início, NOT NULL |
| end_date | TIMESTAMP | Data fim, NOT NULL |
| alert_at | DECIMAL(5,2) | % para alerta (ex: 80.0) |
| is_active | BOOLEAN | Orçamento ativo? |
| description | TEXT | Descrição |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

**Índices**: user_id, category_id, start_date, end_date, is_active

## Migrations

O GORM AutoMigrate cria as tabelas automaticamente no startup:

```go
var entities []interface{}
entities = append(entities, auth.Entities()...)
entities = append(entities, budget.Entities()...)

if err := database.AutoMigrate(entities...); err != nil {
    logger.Fatal().Err(err).Msg("failed to migrate database")
}
```

## Seed Data

### Auth Module

Cria automaticamente:

1. **Permissions**: Permissões básicas para cada recurso
2. **Roles**: admin e user
3. **Root User**: username=root, password=root123
4. Atribui role admin ao root

### Budget Module

Cria categorias padrão para o usuário root:

- **Receitas**: Salário, Freelance, Investimentos
- **Despesas**: Alimentação, Transporte, Moradia, Saúde, Lazer, Educação

## Relacionamentos

```sql
-- User → Accounts (1:N)
accounts.user_id → users.id

-- User → Categories (1:N)
categories.user_id → users.id

-- User → Transactions (1:N)
transactions.user_id → users.id

-- Account → Transactions (1:N)
transactions.account_id → accounts.id

-- Category → Transactions (1:N)
transactions.category_id → categories.id

-- Category → Budgets (1:N)
budgets.category_id → categories.id

-- User → Budgets (1:N)
budgets.user_id → users.id
```

## Consultas Comuns

### Saldo de uma Conta

```sql
SELECT balance FROM accounts WHERE id = ? AND user_id = ?;
```

### Transações do Mês

```sql
SELECT * FROM transactions 
WHERE user_id = ? 
  AND EXTRACT(MONTH FROM date) = ? 
  AND EXTRACT(YEAR FROM date) = ?
ORDER BY date DESC;
```

### Gastos por Categoria

```sql
SELECT c.name, SUM(t.amount) as total
FROM transactions t
JOIN categories c ON t.category_id = c.id
WHERE t.user_id = ? 
  AND t.type = 'expense'
  AND t.status = 'completed'
GROUP BY c.name
ORDER BY total DESC;
```

### Verificar Orçamento

```sql
SELECT 
  b.name,
  b.amount as planned,
  b.spent,
  (b.spent / b.amount * 100) as percent_used
FROM budgets b
WHERE b.user_id = ? 
  AND b.is_active = true
  AND CURRENT_DATE BETWEEN b.start_date AND b.end_date;
```

## Próximos Passos

- Explore os [Módulos](modules.md)
- Veja a [API Reference](../api/auth.md)
- Aprenda sobre [Desenvolvimento](../guides/development.md)
