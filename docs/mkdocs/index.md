# Monetics

Bem-vindo à documentação do **Monetics** - uma API RESTful completa para gestão financeira pessoal.

## 🎯 Visão Geral

O Monetics oferece um sistema robusto para controle de finanças pessoais, permitindo aos usuários gerenciar suas contas, categorizar gastos, registrar transações e planejar orçamentos.

## 🚀 Principais Funcionalidades

### Gestão de Contas
Crie e gerencie múltiplas contas bancárias com acompanhamento de saldo em tempo real.

### Categorização Inteligente
Organize suas receitas e despesas em categorias personalizadas para melhor visualização dos gastos.

### Transações Completas
Registre receitas, despesas e transferências entre contas com status e rastreamento detalhado.

### Orçamentos Planejados
Defina limites de gastos por categoria e período (mensal, trimestral, anual) com alertas configuráveis.

### Relatórios Detalhados
Visualize seus gastos mensais, acompanhe orçamentos e analise seu comportamento financeiro.

### Autenticação Robusta
Sistema completo de autenticação e autorização com JWT, roles e permissions granulares.

## 🏗️ Arquitetura

O projeto segue princípios de **Clean Architecture** com separação clara de responsabilidades:

- **Domain**: Entidades de negócio puras
- **Use Cases**: Lógica de aplicação
- **Adapters**: HTTP handlers e repositories
- **Infrastructure**: Banco de dados, validadores, logger

## 📚 Começando

Para começar a usar o Monetics, consulte nosso [Guia de Instalação](getting-started/installation.md).

## 🔗 Links Úteis

- [Repositório GitHub](https://github.com/alpheres/monetics)
- [Swagger Documentation](http://localhost:8080/swagger/index.html)
- [API Reference](api/auth.md)

## 🤝 Contribuindo

Contribuições são bem-vindas! Consulte nosso [Guia de Desenvolvimento](guides/development.md) para mais informações.
