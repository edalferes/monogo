# Monogo

Monogo é um projeto monolítico modular escrito em Go, focado em boas práticas de desenvolvimento.

## Objetivo

Demonstrar como estruturar uma aplicação backend robusta, escalável e de fácil manutenção.utilizando:

- Estrutura modular (domain, repository, usecase, handler)
- Echo Framework v4 para rotas HTTP
- GORM para persistência de dados
- Documentação automática com Swagger
- Respostas HTTP padronizadas
- Validação de dados com go-playground/validator

## Principais Características

- Projeto monolítico, mas com separação clara de responsabilidades
- Modularização por domínio (exemplo: módulo de usuário)
- Pronto para evoluir para microsserviços, se necessário
- Docker Compose para ambiente de desenvolvimento
- Código comentado e documentado

## Como rodar

```sh
make run
```

## Estrutura de Pastas

- `cmd/` - Entrypoint da aplicação
- `internal/` - Código principal, dividido em módulos/domínios
- `pkg/` - Utilitários e helpers reutilizáveis
- `docs/` - Documentação Swagger gerada automaticamente
