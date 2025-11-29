# API do Módulo Auth

Documentação completa dos endpoints de autenticação e gerenciamento de usuários.

## Autenticação

### Login

Autentica um usuário e retorna um token JWT.

**Endpoint**: `POST /v1/auth/login`

**Request Body**:
```json
{
  "username": "root",
  "password": "root123"
}
```

**Response Success (200)**:
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "root",
      "roles": [
        {
          "id": 1,
          "name": "admin",
          "description": "Administrator role"
        }
      ]
    }
  }
}
```

**Response Error (401)**:
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid credentials"
  }
}
```

## Gerenciamento de Usuários

Todos os endpoints de usuários requerem autenticação via JWT no header:
```
Authorization: Bearer <token>
```

### Criar Usuário

**Endpoint**: `POST /v1/auth/users`

**Permissão**: `users:create`

**Request Body**:
```json
{
  "username": "johndoe",
  "password": "securepass123",
  "role_ids": [2]
}
```

**Response Success (201)**:
```json
{
  "success": true,
  "data": {
    "id": 2,
    "username": "johndoe",
    "roles": [
      {
        "id": 2,
        "name": "user"
      }
    ]
  }
}
```

### Listar Usuários

**Endpoint**: `GET /v1/auth/users`

**Permissão**: `users:read`

**Query Parameters**:
- `page` (int): Número da página (default: 1)
- `limit` (int): Itens por página (default: 10)

**Response Success (200)**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "username": "root",
      "roles": [{"id": 1, "name": "admin"}]
    },
    {
      "id": 2,
      "username": "johndoe",
      "roles": [{"id": 2, "name": "user"}]
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 2
  }
}
```

### Buscar Usuário por ID

**Endpoint**: `GET /v1/auth/users/:id`

**Permissão**: `users:read`

**Response Success (200)**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "root",
    "roles": [
      {
        "id": 1,
        "name": "admin",
        "description": "Administrator role"
      }
    ]
  }
}
```

**Response Error (404)**:
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "User not found"
  }
}
```

### Atualizar Usuário

**Endpoint**: `PUT /v1/auth/users/:id`

**Permissão**: `users:update`

**Request Body**:
```json
{
  "password": "newpassword123",
  "role_ids": [1, 2]
}
```

**Response Success (200)**:
```json
{
  "success": true,
  "data": {
    "id": 2,
    "username": "johndoe",
    "roles": [
      {"id": 1, "name": "admin"},
      {"id": 2, "name": "user"}
    ]
  }
}
```

### Deletar Usuário

**Endpoint**: `DELETE /v1/auth/users/:id`

**Permissão**: `users:delete`

**Response Success (204)**:
No content

## Gerenciamento de Roles

### Listar Roles

**Endpoint**: `GET /v1/auth/roles`

**Permissão**: `roles:read`

**Response Success (200)**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "admin",
      "description": "Administrator role",
      "permissions": [
        {
          "id": 1,
          "name": "users:create",
          "resource": "users",
          "action": "create"
        }
      ]
    }
  ]
}
```

### Criar Role

**Endpoint**: `POST /v1/auth/roles`

**Permissão**: `roles:create`

**Request Body**:
```json
{
  "name": "manager",
  "description": "Manager role",
  "permission_ids": [1, 2, 3]
}
```

## Gerenciamento de Permissions

### Listar Permissions

**Endpoint**: `GET /v1/auth/permissions`

**Permissão**: `permissions:read`

**Response Success (200)**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "users:create",
      "description": "Create users",
      "resource": "users",
      "action": "create"
    },
    {
      "id": 2,
      "name": "users:read",
      "description": "Read users",
      "resource": "users",
      "action": "read"
    }
  ]
}
```

### Criar Permission

**Endpoint**: `POST /v1/auth/permissions`

**Permissão**: `permissions:create`

**Request Body**:
```json
{
  "name": "reports:read",
  "description": "Read financial reports",
  "resource": "reports",
  "action": "read"
}
```

## Códigos de Erro

| Código | Descrição |
|--------|-----------|
| `UNAUTHORIZED` | Credenciais inválidas ou token expirado |
| `FORBIDDEN` | Usuário não tem permissão |
| `NOT_FOUND` | Recurso não encontrado |
| `VALIDATION_ERROR` | Dados inválidos |
| `CONFLICT` | Recurso já existe (ex: username duplicado) |
| `INTERNAL_ERROR` | Erro interno do servidor |

## Permissões Padrão

O sistema vem com as seguintes permissões pré-configuradas:

### Users
- `users:create`
- `users:read`
- `users:update`
- `users:delete`

### Roles
- `roles:create`
- `roles:read`
- `roles:update`
- `roles:delete`

### Permissions
- `permissions:create`
- `permissions:read`
- `permissions:update`
- `permissions:delete`

### Budget (todas as permissões de budget)
- `accounts:*`
- `categories:*`
- `transactions:*`
- `budgets:*`
- `reports:*`

O role **admin** tem todas as permissões.  
O role **user** tem permissões de leitura e escrita no módulo budget apenas.
