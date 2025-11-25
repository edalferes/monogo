# TechDocs - Monetics

Este diretório contém a documentação técnica do Monetics usando MkDocs para integração com o Backstage TechDocs.

## 📚 Estrutura

```
docs/
├── index.md                    # Página inicial
├── getting-started/            # Começando
│   ├── installation.md
│   ├── configuration.md
│   └── quickstart.md
├── architecture/               # Arquitetura
│   ├── overview.md
│   ├── modules.md
│   └── database.md
├── modules/                    # Documentação dos módulos
│   ├── auth.md
│   └── budget.md
├── api/                        # Referência da API
│   ├── auth.md
│   ├── accounts.md
│   ├── categories.md
│   ├── transactions.md
│   ├── budgets.md
│   └── reports.md
└── guides/                     # Guias
    ├── development.md
    ├── testing.md
    └── deployment.md
```

## 🚀 Como Usar Localmente

### Instalar MkDocs

```bash
# Via pip
pip install mkdocs-techdocs-core

# Ou via pipx (recomendado)
pipx install mkdocs
pipx inject mkdocs mkdocs-techdocs-core
```

### Servir Localmente

```bash
# No diretório raiz do projeto
mkdocs serve

# Acesse: http://localhost:8000
```

### Build da Documentação

```bash
mkdocs build

# Gera a pasta site/ com HTML estático
```

## 📖 Backstage Integration

### Configuração no Backstage

A anotação no `catalog-info.yaml` aponta para esta documentação:

```yaml
metadata:
  annotations:
    backstage.io/techdocs-ref: dir:.
```

Isso indica que o TechDocs deve buscar o `mkdocs.yml` na raiz do repositório.

### Como o Backstage Processa

1. **Discovery**: Backstage encontra o `catalog-info.yaml`
2. **Build**: Executa `mkdocs build` no repositório
3. **Publish**: Armazena o site gerado
4. **Serve**: Disponibiliza via interface do Backstage

### Visualizar no Backstage

Após registrar o componente:

1. Acesse o componente no Backstage
2. Clique na aba **"Docs"**
3. A documentação será renderizada

## ✍️ Editando a Documentação

### Adicionar Nova Página

1. Crie o arquivo `.md` em `docs/`
2. Adicione ao `nav` em `mkdocs.yml`:

```yaml
nav:
  - Nova Seção:
      - Título: caminho/para/arquivo.md
```

### Sintaxe Markdown

O MkDocs suporta Markdown estendido com:

- **Admonitions** (alertas):
  ```markdown
  !!! warning "Atenção"
      Conteúdo do alerta
  ```

- **Code Blocks** com syntax highlighting:
  ````markdown
  ```go
  func main() {
      fmt.Println("Hello")
  }
  ```
  ````

- **Tabelas**:
  ```markdown
  | Coluna 1 | Coluna 2 |
  |----------|----------|
  | Valor 1  | Valor 2  |
  ```

- **Links internos**:
  ```markdown
  [Texto](../outro-arquivo.md)
  ```

### Preview em Tempo Real

```bash
# MkDocs auto-reload ao salvar
mkdocs serve
```

## 🎨 Personalização

### Tema

Configurado em `mkdocs.yml`:

```yaml
theme:
  name: material
  palette:
    primary: indigo
    accent: indigo
  features:
    - navigation.tabs
    - navigation.instant
    - search.suggest
```

### Plugins

- **search**: Busca integrada
- **techdocs-core**: Compatibilidade com Backstage

## 📝 Boas Práticas

✅ **Mantenha atualizado**: Documente mudanças importantes  
✅ **Seja claro**: Use exemplos e código sempre que possível  
✅ **Organize bem**: Use a estrutura de pastas lógica  
✅ **Links relativos**: Facilita navegação local e no Backstage  
✅ **Imagens**: Coloque em `docs/assets/`  

## 🔗 Links Úteis

- [MkDocs Documentation](https://www.mkdocs.org/)
- [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/)
- [Backstage TechDocs](https://backstage.io/docs/features/techdocs/)
- [TechDocs Core Plugin](https://github.com/backstage/mkdocs-techdocs-core)

## 📦 Estrutura de Arquivos

```
monetics/
├── mkdocs.yml              # Configuração do MkDocs
├── catalog-info.yaml       # Catálogo do Backstage
├── docs/                   # Documentação
│   ├── index.md
│   ├── getting-started/
│   ├── architecture/
│   ├── modules/
│   ├── api/
│   ├── guides/
│   └── assets/            # Imagens, diagramas
└── site/                  # Gerado pelo build (gitignored)
```

## 🚨 Troubleshooting

### Erro ao buildar

```bash
# Instale dependências
pip install mkdocs-techdocs-core

# Limpe o cache
rm -rf site/
mkdocs build
```

### Links quebrados

- Use caminhos relativos: `../outro-arquivo.md`
- Verifique a estrutura em `mkdocs.yml`

### Não aparece no Backstage

- Verifique a anotação em `catalog-info.yaml`
- Confirme que `mkdocs.yml` está na raiz
- Veja os logs do TechDocs no Backstage
