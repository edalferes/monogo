# Estrutura da Documentação

Este diretório contém toda a documentação do projeto organizada por tipo.

## 📁 Estrutura

```
docs/
├── 📘 mkdocs/                      # MkDocs Documentation (Human-readable)
│   ├── index.md                    # Home page
│   ├── MKDOCS.md                   # MkDocs setup guide
│   ├── module-dependencies.md      # DI system documentation
│   ├── testing-dependencies.md     # Testing guide
│   ├── dependency-graph.md         # Dependency visualization
│   ├── http-architecture.md        # HTTP architecture
│   │
│   ├── getting-started/            # Getting started guides
│   │   ├── installation.md
│   │   ├── quickstart.md
│   │   └── configuration.md
│   │
│   ├── architecture/               # Architecture documentation
│   │   ├── overview.md
│   │   ├── modules.md
│   │   ├── communication.md
│   │   └── database.md
│   │
│   ├── api/                        # API documentation
│   │   ├── auth.md
│   │   ├── accounts.md
│   │   └── health.md
│   │
│   └── guides/                     # Development guides
│       ├── development.md
│       ├── testing.md
│       └── deployment.md
│
├── 🔧 openapi/                     # OpenAPI/Swagger specs
│   ├── docs.go                     # Generated Go docs
│   ├── swagger.json                # OpenAPI JSON spec
│   └── swagger.yaml                # OpenAPI YAML spec
│
└── 🧪 postman/                     # API testing
    └── Monetics.postman_collection.json
```

## 📘 Documentação MkDocs

**Propósito**: Documentação legível para desenvolvedores

**Como usar**:
```bash
# Servir localmente em http://127.0.0.1:8000
mkdocs serve

# Gerar site estático
mkdocs build

# Deploy para GitHub Pages
mkdocs gh-deploy
```

**Tecnologias**:
- MkDocs com tema Material
- Markdown com suporte a diagramas Mermaid
- Suporte a modo escuro
- Busca de texto completo

**Adicionar nova página**:
1. Criar arquivo `.md` no subdiretório apropriado
2. Adicionar à navegação do `mkdocs.yml`
3. Mudanças recarregam automaticamente no servidor dev

## 🔧 OpenAPI/Swagger

**Propósito**: Especificação da API e documentação interativa

**Como gerar**:
```bash
# Gerar a partir de anotações no código
make swagger

# Ou manualmente
swag init -g cmd/api/main.go --parseDependency --parseInternal -o docs/openapi
```

**Acessar Swagger UI**:
```
http://localhost:8080/swagger/index.html
```

**Arquivos**:
- `docs.go` - Pacote Go gerado (auto-gerado, não editar)
- `swagger.json` - Especificação OpenAPI 3.0 em JSON
- `swagger.yaml` - Especificação OpenAPI 3.0 em YAML

**Uso no código**:
```go
import _ "github.com/edalferes/monetics/docs/openapi"
```

## 🧪 Collections do Postman

**Propósito**: Testes de API e fluxos de testes manuais

**Como usar**:
1. Importar `Monetics.postman_collection.json` no Postman
2. Configurar variáveis de ambiente:
   - `BASE_URL` - URL base da API (padrão: `http://localhost:8080`)
   - `AUTH_URL` - URL do serviço Auth (para modo microservices)
   - `BUDGET_URL` - URL do serviço Budget (para modo microservices)
3. Executar collections ou requisições individuais

**Inclui**:
- Fluxos de autenticação (registro, login)
- Gestão de orçamento (contas, categorias, transações)
- Endpoints de health check
- Scripts pré-requisição para gestão de tokens

**Atualizar collection**:
1. Fazer mudanças no Postman
2. Exportar collection (Collection v2.1)
3. Substituir `postman/Monetics.postman_collection.json`

## 🔄 Manutenção

### Atualizando Documentação Swagger

Quando adicionar/modificar endpoints da API:

1. Adicionar anotações Swagger ao handler:
```go
// @Summary Criar conta
// @Description Criar uma nova conta financeira
// @Tags Contas
// @Accept json
// @Produce json
// @Param account body dto.CreateAccountRequest true "Dados da conta"
// @Success 201 {object} dto.AccountResponse
// @Router /v1/budget/accounts [post]
func CreateAccount(c echo.Context) error { ... }
```

2. Regenerar documentação:
```bash
make swagger
```

3. Commitar mudanças:
```bash
git add docs/openapi/
git commit -m "docs: atualizar specs do swagger"
```

### Atualizando MkDocs

Quando adicionar nova documentação:

1. Criar/editar arquivos `.md` no subdiretório apropriado
2. Adicionar à seção nav do `mkdocs.yml`
3. Testar localmente:
```bash
mkdocs serve
```

4. Commitar mudanças:
```bash
git add docs/ mkdocs.yml
git commit -m "docs: adicionar nova documentação"
```

### Atualizando Collection do Postman

Quando a API mudar:

1. Atualizar requisições no Postman
2. Testar todos os endpoints
3. Exportar collection
4. Substituir arquivo em `docs/postman/`
5. Commitar:
```bash
git add docs/postman/
git commit -m "docs: atualizar collection do postman"
```

## 📊 Comparação dos Tipos de Documentação

| Tipo | Formato | Auto-gerado | Público-alvo | Caso de Uso |
|------|--------|---------------|-----------------|----------|
| **MkDocs** | Markdown | ❌ Manual | Desenvolvedores | Arquitetura, guias, tutoriais |
| **Swagger** | YAML/JSON | ✅ Do código | Desenvolvedores, consumidores API | Referência API, testes |
| **Postman** | JSON | ❌ Manual | QA, Desenvolvedores | Testes manuais, testes integração |

## 🎯 Acesso Rápido

- **Docs Locais**: http://127.0.0.1:8000 (executar `mkdocs serve`)
- **Swagger UI**: http://localhost:8080/swagger/index.html (quando API estiver rodando)
- **Postman Collection**: Importar `docs/postman/Monetics.postman_collection.json`

## 📝 Boas Práticas

### MkDocs
- ✅ Usar títulos descritivos
- ✅ Adicionar exemplos de código
- ✅ Incluir diagramas (Mermaid)
- ✅ Manter páginas focadas (um tópico por página)
- ✅ Adicionar breadcrumbs de navegação
- ❌ Não duplicar referência da API (usar link do Swagger)

### Swagger
- ✅ Documentar todos os endpoints
- ✅ Incluir exemplos de request/response
- ✅ Adicionar descrições aos parâmetros
- ✅ Agrupar endpoints relacionados com tags
- ✅ Documentar respostas de erro
- ❌ Não escrever explicações longas (usar MkDocs)

### Postman
- ✅ Organizar em pastas por módulo
- ✅ Usar variáveis de ambiente
- ✅ Adicionar scripts pré-requisição para auth
- ✅ Incluir respostas de exemplo
- ✅ Adicionar asserções de teste
- ❌ Não usar credenciais hardcoded

## 🔗 Links Externos

- [Documentação MkDocs](https://www.mkdocs.org/)
- [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/)
- [Especificação Swagger/OpenAPI](https://swagger.io/specification/)
- [Swag (Go Swagger)](https://github.com/swaggo/swag)
- [Documentação Postman](https://learning.postman.com/docs/)
