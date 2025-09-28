# Configuração - Clean Architecture

Esta implementação segue os princípios de Clean Code e SOLID, separando responsabilidades em módulos específicos.

**Nota**: Esta configuração está em `internal/` porque é específica para este projeto (contém tipos como DatabaseConfig, JWTConfig). Para reutilização em outros projetos, seria necessário criar uma versão mais genérica em `pkg/`.

## Estrutura

```
internal/config/
├── types.go      # Definição dos tipos e structs de configuração
├── loader.go     # Responsável por carregar configurações (SRP)
├── validator.go  # Responsável por validar configurações (SRP)
└── config.go     # Funções públicas e fallbacks
```

## Responsabilidades Separadas

### 1. **Types** (`types.go`)
- ✅ Define apenas as estruturas de dados
- ✅ Métodos utilitários específicos da Config
- ✅ Single Responsibility: Definir tipos

### 2. **Loader** (`loader.go`) 
- ✅ Carrega configurações de arquivos YAML
- ✅ Gerencia variáveis de ambiente
- ✅ Define valores padrão
- ✅ Single Responsibility: Carregar dados

### 3. **Validator** (`validator.go`)
- ✅ Valida se configuração é válida
- ✅ Regras de negócio de validação
- ✅ Single Responsibility: Validar dados

### 4. **Config** (`config.go`)
- ✅ API pública simples
- ✅ Coordena Loader + Validator
- ✅ Fallbacks para compatibilidade

## Princípios SOLID Atendidos

- **SRP**: Cada classe tem uma responsabilidade
- **OCP**: Extensível sem modificar código existente
- **LSP**: Substitution principle respeitado
- **ISP**: Interfaces pequenas e específicas
- **DIP**: Depende de abstrações, não implementações

## Uso

```go
// Básico
config := config.LoadConfig()

// Com opções customizadas
loader := config.NewLoader()
validator := config.NewValidator()

cfg, err := loader.Load(config.LoadOptions{
    ConfigName: "production",
})
if err != nil {
    return err
}

if err := validator.Validate(cfg); err != nil {
    return err
}
```

## Benefícios

1. **Testabilidade**: Cada componente pode ser testado isoladamente
2. **Manutenibilidade**: Mudanças são localizadas
3. **Extensibilidade**: Fácil adicionar novos tipos de config
4. **Clareza**: Cada arquivo tem propósito bem definido