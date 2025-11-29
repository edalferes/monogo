# MkDocs Documentation

Este projeto usa **MkDocs** com o tema **Material** para documentação.

## 📚 Documentação Disponível

- **Module Dependencies** - Sistema de injeção de dependências
- **Testing Dependencies** - Guia de testes de dependências
- **Dependency Graph** - Visualização do grafo de dependências
- **Communication Patterns** - Padrões de comunicação (Local vs HTTP)
- **Health Checks API** - Endpoints de health check para Kubernetes

## 🚀 Visualizar Documentação Localmente

### Pré-requisitos

```bash
# Instalar MkDocs e dependências
pip install mkdocs mkdocs-material pymdown-extensions
```

### Servir Localmente

```bash
# Opção 1: Usar script
./scripts/docs-serve.sh

# Opção 2: Comando direto
mkdocs serve
```

Acesse: http://127.0.0.1:8000

### Build Estático

```bash
# Gerar site estático em ./site
mkdocs build

# Visualizar site gerado
open site/index.html
```

## 📖 Estrutura da Documentação

```
docs/
├── index.md                           # Home page
├── module-dependencies.md             # ✅ Novo - Sistema de DI
├── testing-dependencies.md            # ✅ Novo - Testes
├── dependency-graph.md                # ✅ Novo - Grafo visual
├── getting-started/
│   ├── installation.md
│   ├── quickstart.md
│   └── configuration.md
├── architecture/
│   ├── overview.md
│   ├── modules.md
│   ├── database.md
│   └── communication.md               # ✅ Novo - Padrões de comunicação
├── modules/
│   ├── auth.md
│   └── budget.md
├── api/
│   ├── auth.md
│   ├── accounts.md
│   └── health.md                      # ✅ Novo - Health checks
└── guides/
    ├── development.md
    ├── testing.md
    └── deployment.md
```

## 🎨 Features do Tema Material

- ✅ **Dark Mode**: Alternância entre modo claro/escuro
- ✅ **Search**: Busca instantânea na documentação
- ✅ **Code Highlighting**: Syntax highlight para múltiplas linguagens
- ✅ **Navigation Tabs**: Navegação organizada em tabs
- ✅ **Mermaid Diagrams**: Suporte a diagramas Mermaid
- ✅ **Admonitions**: Blocos de aviso, nota, dica, etc
- ✅ **Mobile Responsive**: Interface adaptável para mobile

## 📝 Adicionando Nova Documentação

### 1. Criar arquivo Markdown

```bash
# Criar nova página
touch docs/new-page.md
```

### 2. Adicionar ao mkdocs.yml

```yaml
nav:
  - Home: index.md
  - Your Section:
      - New Page: new-page.md
```

### 3. Visualizar mudanças

```bash
mkdocs serve
# MkDocs recarrega automaticamente ao salvar arquivos
```

## 🚢 Deploy

### GitHub Pages

```bash
# Deploy para GitHub Pages (branch gh-pages)
mkdocs gh-deploy
```

### Docker

```bash
# Build container com docs
docker build -t monetics-docs -f Dockerfile.docs .

# Run
docker run -p 8000:8000 monetics-docs
```

### Netlify / Vercel

```bash
# Build command
mkdocs build

# Publish directory
site/
```

## 🔗 Links Úteis

- [MkDocs Documentation](https://www.mkdocs.org/)
- [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/)
- [Markdown Guide](https://www.markdownguide.org/)
- [Mermaid Diagrams](https://mermaid.js.org/)

## 📊 Exemplos de Uso

### Admonitions

```markdown
!!! note "Nota"
    Este é um bloco de nota

!!! warning "Atenção"
    Este é um aviso importante

!!! tip "Dica"
    Esta é uma dica útil
```

### Code Blocks

```markdown
​```go
func main() {
    fmt.Println("Hello, World!")
}
​```
```

### Mermaid Diagrams

```markdown
​```mermaid
graph TD
    A[Start] --> B[Process]
    B --> C[End]
​```
```

### Tabs

```markdown
=== "Go"
    ​```go
    fmt.Println("Hello")
    ​```

=== "Python"
    ​```python
    print("Hello")
    ​```
```

## 🐛 Troubleshooting

### "mkdocs: command not found"

```bash
# Instalar MkDocs
pip install mkdocs mkdocs-material
```

### Warnings sobre links quebrados

Verifique se os arquivos referenciados no `mkdocs.yml` existem em `docs/`.

### Site não atualiza automaticamente

```bash
# Reiniciar servidor MkDocs
Ctrl+C
mkdocs serve
```
