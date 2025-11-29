# Guia de Desenvolvimento

Este guia fornece informações para desenvolvedores que desejam contribuir ou estender o Monetics.

## Configurando o Ambiente

### 1. Clone e Configure

```bash
git clone https://github.com/alpheres/monetics.git
cd monetics
go mod download
```

### 2. Configure o Banco de Dados

```bash
# Via Docker
docker-compose up -d postgres

# Ou localmente
createdb monetics
```

### 3. Execute em Modo de Desenvolvimento

```bash
make run
```

## Estrutura do Código

### Adicionando um Novo Use Case

1. Crie o arquivo em `internal/modules/{module}/usecase/`
2. Defina interfaces necessárias em `usecase/interfaces/`
3. Implemente a lógica de negócio
4. Injete via WireUp

Exemplo:

```go
// usecase/create_something.go
package usecase

type CreateSomethingInput struct {
    Name string `validate:"required"`
}

type CreateSomethingUseCase struct {
    repo SomethingRepository
}

func NewCreateSomethingUseCase(repo SomethingRepository) *CreateSomethingUseCase {
    return &CreateSomethingUseCase{repo: repo}
}

func (uc *CreateSomethingUseCase) Execute(input CreateSomethingInput) error {
    // Validação
    // Lógica de negócio
    // Persistência
    return uc.repo.Create(something)
}
```

### Adicionando um Novo Endpoint

1. Crie o handler em `adapters/http/handlers/`
2. Defina DTOs em `adapters/http/dto/`
3. Registre a rota no `module.go`

Exemplo:

```go
// adapters/http/handlers/something_handler.go
package handlers

type SomethingHandler struct {
    useCase *usecase.CreateSomethingUseCase
}

// Create godoc
// @Summary Create something
// @Tags something
// @Accept json
// @Produce json
// @Param input body dto.CreateSomethingDTO true "Input"
// @Success 201 {object} responses.SuccessResponse
// @Router /v1/something [post]
// @Security BearerAuth
func (h *SomethingHandler) Create(c echo.Context) error {
    var dto dto.CreateSomethingDTO
    if err := c.Bind(&dto); err != nil {
        return responses.BadRequest(c, "Invalid input")
    }
    
    result, err := h.useCase.Execute(dto.ToInput())
    if err != nil {
        return responses.InternalError(c, err.Error())
    }
    
    return responses.Created(c, result)
}
```

## Swagger Documentation

### Anotações Swagger

Use comentários GoDoc com anotações Swag:

```go
// @Summary Breve descrição
// @Description Descrição detalhada
// @Tags tag-name
// @Accept json
// @Produce json
// @Param name path string true "Description"
// @Param body body dto.InputDTO true "Body"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /v1/endpoint [post]
// @Security BearerAuth
```

### Gerar Swagger

```bash
make swagger
```

## Testes

### Estrutura de Testes

```
module/
├── usecase/
│   ├── create_something.go
│   └── create_something_test.go
└── adapters/
    └── repository/
        ├── repository.go
        └── repository_test.go
```

### Executar Testes

```bash
# Todos os testes
go test ./...

# Com coverage
go test -cover ./...

# Teste específico
go test ./internal/modules/budget/usecase/
```

### Exemplo de Teste

```go
func TestCreateSomething(t *testing.T) {
    // Arrange
    mockRepo := &MockRepository{}
    useCase := NewCreateSomethingUseCase(mockRepo)
    
    input := CreateSomethingInput{
        Name: "Test",
    }
    
    // Act
    err := useCase.Execute(input)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, 1, mockRepo.CreateCallCount)
}
```

## Migrations

O Monetics usa AutoMigrate do GORM. Para adicionar uma nova entidade:

1. Defina a struct em `domain/`
2. Adicione aos `Entities()` do módulo
3. Reinicie a aplicação

```go
// module.go
func Entities() []interface{} {
    return []interface{}{
        &domain.Something{},
    }
}
```

## Logging

Use o logger injetado:

```go
logger.Info().
    Str("user_id", userID).
    Msg("User created successfully")

logger.Error().
    Err(err).
    Msg("Failed to create user")
```

## Padrões de Código

### Naming Conventions

- **Packages**: lowercase, sem underscores
- **Files**: snake_case
- **Types**: PascalCase
- **Functions**: camelCase (exported) ou camelCase (private)
- **Constants**: PascalCase ou SCREAMING_SNAKE_CASE

### Error Handling

```go
// Retorne erros específicos
if err != nil {
    return errors.New("failed to create user")
}

// Use errors customizados quando necessário
return &errors.ValidationError{
    Field: "username",
    Message: "username already exists",
}
```

### Response Patterns

Use o package `responses`:

```go
responses.Success(c, data)           // 200
responses.Created(c, data)           // 201
responses.NoContent(c)               // 204
responses.BadRequest(c, msg)         // 400
responses.Unauthorized(c, msg)       // 401
responses.Forbidden(c, msg)          // 403
responses.NotFound(c, msg)           // 404
responses.InternalError(c, msg)      // 500
```

## Boas Práticas

✅ Sempre valide inputs  
✅ Use DTOs para separar API de domínio  
✅ Retorne erros específicos  
✅ Log eventos importantes  
✅ Escreva testes unitários  
✅ Documente endpoints com Swagger  
✅ Siga o padrão arquitetural existente  

## Contribuindo

1. Fork o projeto
2. Crie uma branch (`git checkout -b feature/amazing-feature`)
3. Commit suas mudanças (`git commit -m 'Add amazing feature'`)
4. Push para a branch (`git push origin feature/amazing-feature`)
5. Abra um Pull Request

## Próximos Passos

- Veja os [Módulos](../architecture/modules.md)
- Entenda o [Banco de Dados](../architecture/database.md)
- Explore a [API Reference](../api/auth.md)
