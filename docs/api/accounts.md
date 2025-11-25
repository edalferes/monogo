# API de Contas

Gerenciamento de contas bancárias.

## Criar Conta

**Endpoint**: `POST /v1/budget/accounts`

**Headers**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body**:
```json
{
  "name": "Conta Corrente",
  "type": "checking",
  "balance": 5000.00,
  "currency": "BRL",
  "is_active": true,
  "description": "Minha conta principal"
}
```

**Tipos de Conta**:
- `checking` - Conta Corrente
- `savings` - Poupança
- `investment` - Investimento
- `credit_card` - Cartão de Crédito

**Response Success (201)**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "name": "Conta Corrente",
    "type": "checking",
    "balance": 5000.00,
    "currency": "BRL",
    "is_active": true,
    "description": "Minha conta principal",
    "created_at": "2025-11-25T10:00:00Z",
    "updated_at": "2025-11-25T10:00:00Z"
  }
}
```

## Listar Contas

**Endpoint**: `GET /v1/budget/accounts`

**Query Parameters**:
- `is_active` (bool): Filtrar por status ativo

**Response Success (200)**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Conta Corrente",
      "type": "checking",
      "balance": 5000.00,
      "currency": "BRL",
      "is_active": true
    }
  ]
}
```

## Buscar Conta por ID

**Endpoint**: `GET /v1/budget/accounts/:id`

**Response Success (200)**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "name": "Conta Corrente",
    "type": "checking",
    "balance": 5000.00,
    "currency": "BRL",
    "is_active": true,
    "created_at": "2025-11-25T10:00:00Z"
  }
}
```

## Consultar Saldo

**Endpoint**: `GET /v1/budget/accounts/:id/balance`

**Response Success (200)**:
```json
{
  "success": true,
  "data": {
    "account_id": 1,
    "balance": 5000.00,
    "currency": "BRL"
  }
}
```

## Atualizar Conta

**Endpoint**: `PUT /v1/budget/accounts/:id`

**Request Body**:
```json
{
  "name": "Conta Corrente Principal",
  "is_active": true
}
```

## Deletar Conta

**Endpoint**: `DELETE /v1/budget/accounts/:id`

**Response Success (204)**:
No content

!!! warning "Atenção"
    Deletar uma conta não deleta as transações associadas. Considere desativar em vez de deletar.
