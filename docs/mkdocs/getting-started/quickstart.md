# Primeiro Uso

Este guia rápido mostra como fazer suas primeiras chamadas à API do Monetics.

## 1. Autenticação

Primeiro, faça login com o usuário root padrão:

```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "root",
    "password": "root123"
  }'
```

Resposta:

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "root",
      "roles": ["admin"]
    }
  }
}
```

Salve o token JWT retornado para usar nas próximas requisições.

## 2. Criar uma Conta

Crie sua primeira conta bancária:

```bash
curl -X POST http://localhost:8080/v1/budget/accounts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN_AQUI" \
  -d '{
    "name": "Conta Corrente",
    "type": "checking",
    "balance": 5000.00,
    "currency": "BRL",
    "is_active": true
  }'
```

## 3. Criar Categorias

Crie categorias para organizar suas transações:

```bash
# Categoria de Alimentação
curl -X POST http://localhost:8080/v1/budget/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN_AQUI" \
  -d '{
    "name": "Alimentação",
    "type": "expense",
    "color": "#FF5722",
    "icon": "restaurant"
  }'

# Categoria de Salário
curl -X POST http://localhost:8080/v1/budget/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN_AQUI" \
  -d '{
    "name": "Salário",
    "type": "income",
    "color": "#4CAF50",
    "icon": "work"
  }'
```

## 4. Registrar Transações

### Receita (Salário)

```bash
curl -X POST http://localhost:8080/v1/budget/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN_AQUI" \
  -d '{
    "account_id": 1,
    "category_id": 2,
    "type": "income",
    "amount": 5000.00,
    "description": "Salário Novembro",
    "date": "2025-11-25T00:00:00Z",
    "status": "completed"
  }'
```

### Despesa (Supermercado)

```bash
curl -X POST http://localhost:8080/v1/budget/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN_AQUI" \
  -d '{
    "account_id": 1,
    "category_id": 1,
    "type": "expense",
    "amount": 350.00,
    "description": "Supermercado",
    "date": "2025-11-25T00:00:00Z",
    "status": "completed"
  }'
```

## 5. Criar um Orçamento

Defina um limite mensal para alimentação:

```bash
curl -X POST http://localhost:8080/v1/budget/budgets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN_AQUI" \
  -d '{
    "category_id": 1,
    "name": "Orçamento Alimentação Novembro",
    "amount": 1500.00,
    "period": "monthly",
    "start_date": "2025-11-01T00:00:00Z",
    "end_date": "2025-11-30T23:59:59Z",
    "alert_at": 80.0,
    "is_active": true
  }'
```

## 6. Consultar Saldo da Conta

```bash
curl -X GET http://localhost:8080/v1/budget/accounts/1/balance \
  -H "Authorization: Bearer SEU_TOKEN_AQUI"
```

## 7. Ver Relatório Mensal

Obtenha um resumo dos seus gastos:

```bash
curl -X GET "http://localhost:8080/v1/budget/reports/monthly?month=11&year=2025" \
  -H "Authorization: Bearer SEU_TOKEN_AQUI"
```

## Swagger UI

Para uma experiência interativa, acesse o Swagger UI:

```
http://localhost:8080/swagger/index.html
```

Lá você pode:

- Ver todos os endpoints disponíveis
- Testar as APIs diretamente no navegador
- Ver exemplos de request/response
- Entender os modelos de dados

## Próximos Passos

- Explore todos os endpoints na [API Reference](../api/auth.md)
- Entenda a arquitetura em [Visão Geral](../architecture/overview.md)
- Configure orçamentos avançados consultando a documentação da API
